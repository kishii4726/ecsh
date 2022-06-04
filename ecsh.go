package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/manifoldco/promptui"
)

func chooseValueFromItems(label string, items []string) string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, selected_value, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return selected_value
}

func getEcsClusters(client *ecs.Client) []string {
	list_clusters, err := client.ListClusters(context.TODO(), &ecs.ListClustersInput{})
	if err != nil {
		log.Fatalf("ListClusters failed %v\n", err)
	}
	ecs_cluster_arns := list_clusters.ClusterArns
	if len(ecs_cluster_arns) == 0 {
		log.Fatalf("Cluster does not exist")
	}
	ecs_clusters := []string{}
	for _, v := range ecs_cluster_arns {
		ecs_clusters = append(ecs_clusters, strings.Split(v, "/")[1])
	}

	return ecs_clusters
}

func getEcsServices(client *ecs.Client, ecs_cluster string) []string {
	list_services, err := client.ListServices(context.TODO(), &ecs.ListServicesInput{
		Cluster: aws.String(ecs_cluster),
	})
	if err != nil {
		log.Fatalf("ListServices failed %v\n", err)
	}
	ecs_service_arns := list_services.ServiceArns
	if len(ecs_service_arns) == 0 {
		log.Fatalf("Service does not exist")
	}
	ecs_services := []string{}
	for _, v := range ecs_service_arns {
		ecs_services = append(ecs_services, strings.Split(v, "/")[2])
	}

	return ecs_services
}

func getEcsTaskIds(client *ecs.Client, ecs_cluster string, ecs_service string) []string {
	list_tasks, err := client.ListTasks(context.TODO(), &ecs.ListTasksInput{
		Cluster:     aws.String(ecs_cluster),
		ServiceName: aws.String(ecs_service),
	})
	if err != nil {
		log.Fatalf("ListTasks failed %v\n", err)
	}
	ecs_task_arns := list_tasks.TaskArns
	if len(ecs_task_arns) == 0 {
		log.Fatalf("Task does not exist")
	}
	ecs_task_ids := []string{}
	for _, v := range ecs_task_arns {
		ecs_task_ids = append(ecs_task_ids, strings.Split(v, "/")[2])
	}

	return ecs_task_ids
}

func main() {

	// select region
	aws_region := chooseValueFromItems("Select Region", []string{"ap-northeast-1", "us-east-1", "us-east-2"})

	// set config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(aws_region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := ecs.NewFromConfig(cfg)

	// select cluster
	ecs_cluster := chooseValueFromItems("Select ECS Cluster", getEcsClusters(client))
	ecs_service := chooseValueFromItems("Select ECS Service", getEcsServices(client, ecs_cluster))
	ecs_task_id := chooseValueFromItems("Select ECS Task Id", getEcsTaskIds(client, ecs_cluster, ecs_service))

	// select container
	describe_tasks, err := client.DescribeTasks(context.TODO(), &ecs.DescribeTasksInput{
		Tasks:   []string{ecs_task_id},
		Cluster: aws.String(ecs_cluster),
	})
	ecs_containers := describe_tasks.Tasks[0].Containers
	fmt.Printf("%T\n", ecs_containers)
	ecs_container_names := []string{}
	for _, v := range ecs_containers {
		ecs_container_names = append(ecs_container_names, *v.Name)
	}
	ecs_container := chooseValueFromItems("Select ECS Container", ecs_container_names)

	// get runtimeId
	var ecs_runtime_id string
	for _, v := range ecs_containers {
		if *v.Name == ecs_container {
			ecs_runtime_id = strings.Split(*v.RuntimeId, "-")[0]
		}
	}

	// select shell
	ecs_shell := chooseValueFromItems("Select Shell", []string{"sh", "bash"})

	// execute command
	out, err := client.ExecuteCommand(context.TODO(), &ecs.ExecuteCommandInput{
		Command:     aws.String(ecs_shell),
		Interactive: true,
		Task:        aws.String(ecs_task_id),
		Cluster:     aws.String(ecs_cluster),
		Container:   aws.String(ecs_container),
	})
	sess, _ := json.Marshal(out.Session)
	var target = fmt.Sprintf("ecs:%s_%s_%s", ecs_cluster, ecs_task_id, ecs_runtime_id)
	var ssmTarget = ssm.StartSessionInput{
		Target: &target,
	}
	targetJSON, err := json.Marshal(ssmTarget)

	// start session
	cmd := exec.Command(
		"session-manager-plugin",
		string(sess),
		aws_region,
		"StartSession",
		"",
		string(targetJSON),
		"https://ssm."+aws_region+".amazonaws.com",
	)
	signal.Ignore(os.Interrupt)
	defer signal.Reset(os.Interrupt)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Run()
}