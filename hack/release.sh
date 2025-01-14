#!/usr/bin/env bash

# Copyright 2018 The Knative Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

source $(dirname $0)/../vendor/github.com/knative/test-infra/scripts/release.sh

# Set default GCS/GCR
: ${EVENTING_SOURCES_RELEASE_GCS:="knative-releases/eventing-sources"}
: ${EVENTING_SOURCES_RELEASE_GCR:="gcr.io/knative-releases"}
readonly EVENTING_SOURCES_RELEASE_GCS
readonly EVENTING_SOURCES_RELEASE_GCR

# TODO(n3wscott): finalize the release packages when we have something to release.

# Yaml files to generate, and the source config dir for them.
declare -A RELEASES
RELEASES["release.yaml"]="config"
readonly RELEASES

# Script entry point.

parse_flags $@

set -o errexit
set -o pipefail

run_validation_tests ./test/presubmit-tests.sh

banner "Building the release"

# Set the repository
export KO_DOCKER_REPO=${EVENTING_SOURCES_RELEASE_GCR}

if (( PUBLISH_RELEASE )); then
  echo "- Destination GCR: ${EVENTING_SOURCES_RELEASE_GCR}"
  echo "- Destination GCS: ${EVENTING_SOURCES_RELEASE_GCS}"
fi

# Build the release

all_yamls=()

for yaml in "${!RELEASES[@]}"; do
  config="${RELEASES[${yaml}]}"
  echo "Building Knative Eventing Sources - ${config}"
  ko resolve ${KO_FLAGS} -f ${config}/ > ${yaml}
  tag_images_in_yaml ${yaml} ${EVENTING_SOURCES_RELEASE_GCR} ${TAG}
  all_yamls+=(${yaml})
done

echo "New release built successfully"

if (( ! PUBLISH_RELEASE )); then
 exit 0
fi

# Publish the release

for yaml in ${all_yamls[@]}; do
  echo "Publishing ${yaml}"
  publish_yaml ${yaml} ${EVENTING_SOURCES_RELEASE_GCS} ${TAG}
done

echo "New release published successfully"
