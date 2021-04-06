package utils

import (
        "context"
        "fmt"
        "grafana/pkg/notification"
        "runtime"
        "time"
        //  "github.com/aws/aws-sdk-go-v2/aws"
        "github.com/aws/aws-sdk-go-v2/config"
        "github.com/aws/aws-sdk-go-v2/service/ec2"
)

func init() {
        UTILS = append(UTILS, EventCheck)
}

func EventCheck(ctx context.Context) {
        info.Println("Start ec2 instance event check.")
        for {
                checked := Event()
                for _, v := range checked {
                        if v != "" {
                                if err := notification.Text(v); err != nil {
                                        info.Println("dingding send error: ", err)
                                }
                        }
                }
                time.Sleep(24 * time.Hour)
                select {
                case <-ctx.Done():
                        info.Println("done")
                        return
                default:
                }
        }
}

func Event() []string {
        events := []string{}
        cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
        if err != nil {
                info.Fatalf("unable to load SDK config, %v", err)
        }
        client := ec2.NewFromConfig(cfg)
        input := &ec2.DescribeInstanceStatusInput{}
        data, err := client.DescribeInstanceStatus(context.Background(), input)
        if err != nil {
                info.Println(err)
        }

        for _, instancestatus := range data.InstanceStatuses {
                instanceid := *instancestatus.InstanceId
                input := &ec2.DescribeInstancesInput{InstanceIds: []string{instanceid}}
                instancesdata, err := client.DescribeInstances(context.Background(), input)
                if err != nil {
                        info.Println(err)
                }
                var (
                        ip          string
                        description string
                        timeafter   time.Time
                        name        string
                )

                for _, i := range instancesdata.Reservations {
                        for _, instance := range i.Instances {
                                ip = *instance.PrivateIpAddress + " " + ip
                                for _, tag := range instance.Tags {
                                        if *tag.Key == "Name" {
                                                name = *tag.Value
                                        }
                                }
                        }
                }

                for _, instancevent := range instancestatus.Events {
                        description = *(instancevent.Description) + " " + description
                        timeafter = *(instancevent.NotAfter)
                }
                msg := `EC2 Events %s alert
Role: %s
InstanceId: %s
IP: %s
ScheduleTime: %s
`
                if description != "" {
                        events = append(events, fmt.Sprintf(msg, description, name, instanceid, ip, timeafter.String()))

                }
        }
        return events

}
