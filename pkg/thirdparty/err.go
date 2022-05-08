package thirdparty

import "git.cestong.com.cn/cecf/cecf-golib/pkg/errors"

var ErrThirdParty = errors.ErrorWithCodeAndHTTPStatus(3598, "third party service failed", 598)
