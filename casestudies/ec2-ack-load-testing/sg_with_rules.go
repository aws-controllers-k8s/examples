package main

import (
	"context"
	"fmt"
	"math/rand"
	"path/filepath"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	totalGroups = 90
	batchSize  = 100
	qps        = 100
	burst      = 100
)

var (
	// Create a new random source seeded with the current time
	randSource = rand.NewSource(time.Now().UnixNano())
	// Create a new random generator
	randGen = rand.New(randSource)
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err.Error())
	}

	config.QPS = float32(qps)
	config.Burst = burst

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	gvr := schema.GroupVersionResource{
		Group:    "ec2.services.k8s.aws",
		Version:  "v1alpha1",
		Resource: "securitygroups",
	}

	startTime := time.Now()

	for batchStart := 0; batchStart < totalGroups; batchStart += batchSize {
		batchEnd := batchStart + batchSize
		if batchEnd > totalGroups {
			batchEnd = totalGroups
		}

		fmt.Printf("creating batch %d to %d\n", batchStart, batchEnd-1)

		// Create a wait group for this batch
		var wg sync.WaitGroup
		wg.Add(batchEnd - batchStart)

		// Create Security Groups in this batch
		for i := batchStart; i < batchEnd; i++ {
			go func(index int) {
				defer wg.Done()

				// Generate a Security Group CR
				group := generateSecurityGroupCR(index)

				// Create the Security Group CR
				_, err := dynamicClient.Resource(gvr).Namespace("default").Create(context.TODO(), group, metav1.CreateOptions{})
				if err != nil {
					fmt.Printf("Error creating Security Group %d: %v\n", index, err)
				} else {
					fmt.Printf("Created Security Group %d\n", index)
				}
			}(i)
		}

		// Wait for this batch to finish
		wg.Wait()
		fmt.Printf("Batch %d to %d completed\n", batchStart, batchEnd-1)
	}
	duration := time.Since(startTime)
	fmt.Printf("all Security Groups have been created in %v \n", duration)

}

func generateSecurityGroupCR(index int) *unstructured.Unstructured {
	sgName := fmt.Sprintf("ack-security-groups-upd-sg-%d", index)

	allPorts := []int{22, 80, 443, 3306, 5432, 6379, 8080, 8443, 1024, 1433, 2049, 2375, 2376, 3389, 5000, 5601, 7000, 7001, 9000, 9090, 9200, 9300, 27017, 28015}

	// Shuffle the ports
	randGen.Shuffle(len(allPorts), func(i, j int) {
		allPorts[i], allPorts[j] = allPorts[j], allPorts[i]
	})

	// Generate 20 ingress rules
	ingressRules := make([]map[string]interface{}, 20)
	for i := 0; i < 20; i++ {
		port := allPorts[i]
		ingressRules[i] = map[string]interface{}{
			"ipProtocol": "tcp",
			"fromPort":   port,
			"toPort":     port,
			"ipRanges": []map[string]interface{}{
				{
					"cidrIP": fmt.Sprintf("%d.%d.%d.%d/32", randGen.Intn(256), randGen.Intn(256), randGen.Intn(256), randGen.Intn(256)),
				},
			},
		}
	}

	// Shuffle the ports again for egress rules
	randGen.Shuffle(len(allPorts), func(i, j int) {
		allPorts[i], allPorts[j] = allPorts[j], allPorts[i]
	})

	// Generate 20 egress rules
	egressRules := make([]map[string]interface{}, 20)
	for i := 0; i < 20; i++ {
		port := allPorts[i]
		egressRules[i] = map[string]interface{}{
			"ipProtocol": "tcp",
			"fromPort":   port,
			"toPort":     port,
			"ipRanges": []map[string]interface{}{
				{
					"cidrIP": "0.0.0.0/0",
				},
			},
		}
	}

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "ec2.services.k8s.aws/v1alpha1",
			"kind":       "SecurityGroup",
			"metadata": map[string]interface{}{
				"name":      sgName,
				"namespace": "default",
			},
			"spec": map[string]interface{}{
				"name":        sgName,
				"vpcID":       "vpc-02780ddde6f72da2b",
				"description": "Security group managed by ACK",
				"tags": []map[string]interface{}{
					{
						"key":   "Name",
						"value": sgName,
					},
				},
				"ingressRules": ingressRules,
				"egressRules":  egressRules,
			},
		},
	}
}
