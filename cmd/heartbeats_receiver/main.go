/*
Copyright 2018 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"github.com/knative/eventing/pkg/event"
	"log"
	"net/http"
	"time"
)

type Heartbeat struct {
	Sequence int    `json:"id"`
	Label    string `json:"label"`
}

func handler(ctx context.Context, hb *Heartbeat) {
	metadata := event.FromContext(ctx)
	log.Printf("[%s] %s %s: %d,%q", metadata.EventTime.Format(time.RFC3339), metadata.ContentType, metadata.Source, hb.Sequence, hb.Label)
}

func main() {
	log.Print("Ready and listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", event.Handler(handler)))
}
