// Copyright 2016 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

/*
Package facility defines the high-level features that Granitic makes available to applications.

A facility is Granitic's term for a group of components that together provide a high-level feature to application developers,
like logging or service error management. This package contains several sub-packages, one for each facility that can be
enabled and configured by user applications.

A full description of how facilities can be enabled and configured can be found at http://granitic.io/1.0/ref/facilities but a basic description
of how they work follows:

Enabling and disabling facilities

The features that are available to applications, and whether they are enabled by default or not, are enumerated in the file:

	$GRANITIC_HOME/resource/facility-config/facilities.json

which will look something like:

	{
	  "Facilities": {
		"HttpServer": false,
		"JsonWs": false,
		"XmlWs": false,
		"FrameworkLogging": true,
		"ApplicationLogging": true,
		"QueryManager": false,
		"RdbmsAccess": false,
		"ServiceErrorManager": false,
		"RuntimeCtl": false,
		"TaskScheduler": false
	  }
	}

This shows that the ApplicationLogging and FrameworkLogging facilities are enabled by default, but everything else is turned
off. If you wanted to enable the HttpServer facility, you'd add the following to any of your application's configuration files:

	{
	  "Facilities": {
		"HttpServer": true
	  }
	}

Configuring facilities

Each facility has a number of default settings that can be found in the file:

	$GRANITIC_HOME/resource/facility-config/facilityname.json

For example, the default configuration for the HttpServer facility will look something like:

  {
    "HttpServer":{
      "Port": 8080,
	  "AccessLogging": false,
	  "TooBusyStatus": 503,
	  "AutoFindHandlers": true
    }
  }

Any of these settings can be changed by overriding one or more of the fields in your application's configuration file. For example, to
change the port on which your application's HTTP server listens on, you could add the following to any of your application's configuration files:

  {
    "HttpServer":{
      "Port": 9000
    }
  }

*/
package facility

import (
	"github.com/graniticio/granitic/config"
	"github.com/graniticio/granitic/ioc"
	"github.com/graniticio/granitic/logging"
)

// A facility builder is responsible for programmatically constructing the objects required to support a facility and,
// where required, adding them to the IoC container.
type FacilityBuilder interface {
	//BuildAndRegister constructs the components that together constitute the facility and stores them in the IoC
	// container.
	BuildAndRegister(lm *logging.ComponentLoggerManager, ca *config.ConfigAccessor, cn *ioc.ComponentContainer) error

	//FacilityName returns the facility's unique name. Used to check whether the facility is enabled in configuration.
	FacilityName() string

	//DependsOnFacilities returns the names of other facilities that must be enabled in order for this facility to run
	//correctly.
	DependsOnFacilities() []string
}
