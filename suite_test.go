package main

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

	
const timeout = 1 * time.Second

func TestPipelines(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Go Pipelines Suite")
}
