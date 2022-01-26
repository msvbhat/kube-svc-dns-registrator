package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func route53CreateRecord(zoneID string, name string, ips []string) bool {
	r53 := route53.New(session.New())
	changes := make([]*route53.Change, 0, len(ips))
	for i, ip := range ips {
		changes = append(changes, &route53.Change{
			Action: aws.String("UPSERT"),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name:            aws.String(name),
				ResourceRecords: []*route53.ResourceRecord{{Value: aws.String(ip)}},
				SetIdentifier:   aws.String(fmt.Sprintf("The name %s index %d", name, i)),
				TTL:             aws.Int64(30),
				Type:            aws.String("A"),
				Weight:          aws.Int64(100),
			},
		})
	}

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
			Comment: aws.String(fmt.Sprintf("The IP Address for service %s", name)),
		},
		HostedZoneId: aws.String(zoneID),
	}

	res, err := r53.ChangeResourceRecordSets(input)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	log.Println("The result is")
	log.Println(res)
	return true
}
