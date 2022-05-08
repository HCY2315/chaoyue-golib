package thirdparty

import "github.com/HCY2315/chaoyue-golib/pkg/errors"

var ErrThirdParty = errors.ErrorWithCodeAndHTTPStatus(3598, "third party service failed", 598)
