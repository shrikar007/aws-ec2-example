package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"os"
)
func init() {
	viper.SetConfigType("toml")
	viper.SetConfigName("config") // name of config file (without extension)

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	viper.AddConfigPath(path)
	err = viper.ReadInConfig() // Find and read the config file
	if err != nil {
		log.Fatal(err)
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		err = viper.ReadInConfig() // Find and read the config file
		if err != nil {
			log.Fatal(err)
		}
	})

	viper.WatchConfig()
}
func main() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials:credentials.NewStaticCredentials(viper.GetString("cred.accesskeyid"),viper.GetString("cred.secretaccesskey"),"")},
	)
	svc := ec2.New(sess)
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:      aws.String("ami-e7527ed7"),
		InstanceType: aws.String("t2.micro"),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
	})
	if err != nil {
		fmt.Println("Could not create instance", err)
		return
	}
	fmt.Println("Created instance", *runResult.Instances[0].InstanceId)
	_, errtag := svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{runResult.Instances[0].InstanceId},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("name"),
				Value: aws.String("amazon"),
			},
		},
	})
	if errtag != nil {
		log.Println("Could not create tags for instance", runResult.Instances[0].InstanceId, errtag)
		return
	}
	fmt.Println("Successfully tagged instance")
}
