// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package container

import (
	"fmt"
	"strings"
)

const DefaultContainerRegistry = "docker.elastic.co"
const DefaultNsApm = "apm"
const DefaultNsElasticsearch = "elasticsearch"
const DefaultNsKibana = "kibana"
const DefaultNsEnterpiseSearch = "enterprise-search"
const DefaultNsBeats = "beats"
const DefaultNsMaps = "elastic-maps-service"

var (
	containerRegistry = DefaultContainerRegistry
	containerSuffix   = ""

	nsApm             = ""
	nsElasticsearch   = ""
	nsKibana          = ""
	nsEnterpiseSearch = ""
	nsBeats           = ""
	nsMaps            = ""
)

// SetContainerRegistry sets the global container registry used to download Elastic stack images.
func SetContainerRegistry(registry string) {
	containerRegistry = registry
}

func SetNsApm(val string) {
	nsApm = val
}

func SetNsElasticsearch(val string) {
	nsElasticsearch = val
}

func SetNsKibana(val string) {
	nsKibana = val
}

func SetNsEnterpiseSearch(val string) {
	nsEnterpiseSearch = val
}

func SetNsBeats(val string) {
	nsBeats = val
}

func SetNsMaps(val string) {
	nsMaps = val
}

func SetContainerSuffix(suffix string) {
	containerSuffix = suffix
}

type Image string

var (
	APMServerImage        Image
	ElasticsearchImage    Image
	KibanaImage           Image
	EnterpriseSearchImage Image
	FilebeatImage         Image
	MetricbeatImage       Image
	HeartbeatImage        Image
	AuditbeatImage        Image
	JournalbeatImage      Image
	PacketbeatImage       Image
	AgentImage            Image
	MapsImage             Image
)

func MakeImageString() {
	APMServerImage = Image(fmt.Sprintf("%s/apm-server", nsApm))
	ElasticsearchImage = Image(fmt.Sprintf("%s/elasticsearch", nsElasticsearch))
	KibanaImage = Image(fmt.Sprintf("%s/kibana", nsKibana))
	EnterpriseSearchImage = Image(fmt.Sprintf("%s/enterprise-search", nsEnterpiseSearch))
	FilebeatImage = Image(fmt.Sprintf("%s/filebeat", nsBeats))
	MetricbeatImage = Image(fmt.Sprintf("%s/metricbeat", nsBeats))
	HeartbeatImage = Image(fmt.Sprintf("%s/heartbeat", nsBeats))
	AuditbeatImage = Image(fmt.Sprintf("%s/auditbeat", nsBeats))
	JournalbeatImage = Image(fmt.Sprintf("%s/journalbeat", nsBeats))
	PacketbeatImage = Image(fmt.Sprintf("%s/packetbeat", nsBeats))
	AgentImage = Image(fmt.Sprintf("%s/elastic-agent", nsBeats))
	MapsImage = Image(fmt.Sprintf("%s/elastic-maps-server-ubi8", nsMaps))
}

// ImageRepository returns the full container image name by concatenating the current container registry and the image path with the given version.
func ImageRepository(img Image, version string) string {
	// don't double append suffix if already contained as e.g. the case for maps
	if strings.HasSuffix(string(img), containerSuffix) {
		return fmt.Sprintf("%s/%s:%s", containerRegistry, img, version)
	}
	return fmt.Sprintf("%s/%s%s:%s", containerRegistry, img, containerSuffix, version)
}
