package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	cfg "github.com/andhikagama/awssh/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
)

var (
	awsAccessID  string
	awsSecretKey string
	awsRegion    string
	tagName      string
	sshUser      string
	awsPem       string
	config       cfg.Config
)

func init() {
	config = cfg.NewViperConfig()
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if config.GetBool(`debug`) {
		logrus.Warn(`ezssh in debug mode`)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	awsAccessID = config.GetString(`aws.access_id`)
	awsSecretKey = config.GetString(`aws.secret_key`)
	awsRegion = config.GetString(`aws.region`)
	tagName = config.GetString(`aws.tag_name`)
	awsPem = config.GetString(`aws.pem`)
	sshUser = config.GetString(`ssh_user`)
}

func main() {
	app := cli.NewApp()
	app.Name = `awssh`
	app.Usage = `SSH to a EC2 instance in AWS without the IP, just the instance:{custom} tag`
	app.HelpName = app.Name
	app.HideHelp = true
	app.HideVersion = true
	app.ArgsUsage = `[instance tag] [key mode]`

	app.Action = func(c *cli.Context) error {
		tag := c.Args().Get(0)
		if tag == `` {
			return cli.ShowAppHelp(c)
		}

		mode := c.Args().Get(1)

		var err error
		if err != nil {
			cli.NewExitError(err.Error(), 1)
		}

		instances, err := getInstances(tag)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if len(instances) == 1 {
			ip := instances[0].PublicIpAddress
			return launchTerminal(*ip, mode)
		}

		servers := map[string]string{}
		options := []string{}
		for _, instance := range instances {
			name := fmt.Sprintf(`%v (%v) `, *instance.InstanceId, *instance.PublicIpAddress)

			for _, tag := range instance.Tags {
				if *tag.Key == `Name` {
					name += *tag.Value
					break
				}
			}

			servers[name] = *instance.PublicIpAddress
			options = append(options, name)
		}

		choice := ``
		prompt := &survey.Select{
			Message: `Multiple instances found, please select one: `,
			Options: options,
		}
		survey.AskOne(prompt, &choice, nil)

		ip := servers[choice]
		return launchTerminal(ip, mode)
	}

	app.Run(os.Args)
}

func getInstances(instanceTag string) ([]*ec2.Instance, error) {
	sess := session.Must(initAWSSession())
	ec2Client := ec2.New(sess)

	filterName := `tag:` + tagName
	filter := ec2.Filter{
		Name: &filterName,
		Values: []*string{
			&instanceTag,
		},
	}

	input := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{&filter},
	}

	result := []*ec2.Instance{}

	output, err := ec2Client.DescribeInstances(&input)
	if err != nil {
		return result, err
	}

	if len(output.Reservations) == 0 {
		return result, errors.New(`No instance found`)
	}

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			if *instance.State.Name == `running` {
				result = append(result, instance)
			}
		}
	}

	if len(result) == 0 {
		return result, errors.New(`No active instance found`)
	}

	return result, nil
}

func initAWSSession() (*session.Session, error) {
	cred := credentials.NewStaticCredentials(awsAccessID, awsSecretKey, ``)
	conf := aws.NewConfig()
	conf.Credentials = cred
	conf.Region = &awsRegion

	return session.NewSession(conf)
}

func launchTerminal(ip string, mode string) error {
	path := fmt.Sprintf(`/tmp/awssh_%v.command`, time.Now().Unix())
	sshCmd := []byte(fmt.Sprintf(`ssh %v@%v`, sshUser, ip))

	if mode == `--key` {
		sshCmd = []byte(fmt.Sprintf(`ssh -i "%v" %v@%v`, awsPem, sshUser, ip))
	}

	err := ioutil.WriteFile(path, sshCmd, 0744)
	if err != nil {
		return err
	}

	cmd := exec.Command(`open`, path)
	return cmd.Run()
}
