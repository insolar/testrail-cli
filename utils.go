// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package testrail_cli

import (
	"path"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func TicketFromURL(url string) string {
	if strings.HasPrefix(url, "https") || strings.HasPrefix(url, "http") {
		s := strings.Split(url, "/")
		return s[len(s)-1]
	}
	return url
}

func TRTicket(id int) string {
	return path.Join(viper.GetString("URL"), "/index.php?/cases/view/", strconv.Itoa(id))
}
