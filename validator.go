package array

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidateOption struct {
	Rule    ValidateRule
	Regular string
	Error   string
}
type ValidateRule = string

const (
	ValidateEmail    ValidateRule = `email`
	ValidatePhone    ValidateRule = `mobile`
	ValidateIp       ValidateRule = `ip`
	ValidateTel      ValidateRule = `tel`
	ValidateUserName ValidateRule = `username`
)

var (
	ValidateInvalidError     = errors.New(`InvalidString`)
	ValidateUnknownRuleError = errors.New(`UnknownRule`)
	ValidateTelError         = errors.New(`TelError`)
	ValidateMobileError      = errors.New(`MobileError`)
	ValidateUserNameError    = errors.New(`UserNameError`)
	ValidateIPError          = errors.New(`IPError`)
	ValidateEmailError       = errors.New(`EmailError`)
	ValidateLengthError      = errors.New(`LengthError`)
	ValidateOptionError      = errors.New(`OptionError`)
	ValidateRegularError     = errors.New(`RegularError`)
)

// Default Rule:email/mobile/ip/tel/username
// Example 1:
// type User struct {
// 	  Name  string `rule:"username,length=1-8,option=XXX YYY Alice"`
//	  Email string `rule:"email"`
//	  Test string `regular:"[0-9]+" error:"MustNumber"`
// }
// username := User{Name: "Alice", Email: "alice@example.com"}
// ok, err := Validate(username) // true

// Example 2:
// ok, err := Validate(`Alice`, ValidateOption{Rule: ValidateUserName})

// Validate support string,slice,struct
func Validate[T interface{}](data T, opts ...ValidateOption) (bool, error) {
	var opt = ValidateOption{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	dataType := reflect.TypeOf(data)
	dataValue := reflect.ValueOf(data)
	switch dataType.Kind() {
	case reflect.String:
		if opt.Rule == `` {
			return validateString(dataValue.String())
		}
		if opt.Regular != `` {
			matchOk := regexp.MustCompile(opt.Regular).MatchString(dataValue.String())
			if !matchOk {
				err := ValidateRegularError
				if opt.Error != `` {
					err = errors.New(opt.Error)
				}
				return false, err
			}
			return matchOk, nil
		}
		return validateStringByRule(dataValue.String(), opt.Rule)
	case reflect.Struct:
		for i := 0; i < dataType.NumField(); i++ {
			field := dataType.Field(i)
			value := dataValue.Field(i)
			rule := field.Tag.Get("rule")
			regular := field.Tag.Get("regular")
			errorTag := field.Tag.Get("error")
			if value.Type().Kind() != reflect.String {
				continue
			}
			if rule != `` && regular == `` {
				result, err := validateStringByRule(value.String(), rule)
				if err != nil || !result {
					if errorTag != `` {
						err = errors.New(errorTag)
					}
					return false, fmt.Errorf("%s:%w", field.Name, err)
				}
			} else if regular != `` {
				matchOk := regexp.MustCompile(regular).MatchString(value.String())
				if !matchOk {
					err := ValidateRegularError
					if errorTag != `` {
						err = errors.New(errorTag)
					}
					return false, err
				}
				return matchOk, nil
			} else {
				return false, errors.New(field.Name + `NotHaveRule`)
			}
		}
		return true, nil
	case reflect.Slice:
		for i := 0; i < dataValue.Len(); i++ {
			element := dataValue.Index(i)
			result, err := Validate(element.Interface())
			if err != nil || !result {
				return false, ValidateInvalidError
			}
		}
		return true, nil
	case reflect.Pointer:
		return Validate(dataValue.Elem().Interface(), opts...)
	default:
		return false, ValidateUnknownRuleError
	}
}

var (
	emailRegexp    = regexp.MustCompile(`^\w+@\w+\.\w+$`)
	ipRegexp       = regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
	usernameRegexp = regexp.MustCompile(`^\w*$`)
	mobileRegexp   = regexp.MustCompile(`^1[3-9]\d{9}$`)
	telRegexp      = regexp.MustCompile(`^\d{3,4}-\d{7,8}$`)
)

func validateStringByRule(str string, rule string) (bool, error) {
	rules := strings.Split(rule, ",")
	for _, r := range rules {
		if sub := regexp.MustCompile(`length\s*=(.*)`).FindStringSubmatch(r); len(sub) > 1 {
			var lengthStr string
			lengthStr = sub[1]
			lengths := strings.Split(lengthStr, "-")
			minLen, _ := strconv.Atoi(lengths[0])
			maxLen := minLen
			if len(lengths) > 1 {
				maxLen, _ = strconv.Atoi(lengths[1])
			}
			if len(str) < minLen || len(str) > maxLen {
				return false, fmt.Errorf(`LengthRange:%s:%w`, lengthStr, ValidateLengthError)
			}
			continue
		}
		if sub := regexp.MustCompile(`option\s*=(.*)`).FindStringSubmatch(r); len(sub) > 1 {
			matchList := strings.Split(sub[0], ` `)
			hasOne := false
			for _, s := range matchList {
				if s == str {
					hasOne = true
					break
				}
			}
			if hasOne == false {
				return false, fmt.Errorf(`OnlyAllow:%s:%w`, sub[0], ValidateOptionError)
			}
			continue
		}
		switch r {
		case "email":
			if !emailRegexp.MatchString(str) {
				return false, ValidateEmailError
			}
		case "ip":
			if !ipRegexp.MatchString(str) {
				return false, ValidateIPError
			}
		case "username":
			if !usernameRegexp.MatchString(str) {
				return false, ValidateUserNameError
			}
		case "mobile":
			if !mobileRegexp.MatchString(str) {
				return false, ValidateMobileError
			}
		case "tel":
			if !telRegexp.MatchString(str) {
				return false, ValidateTelError
			}
		default:
			return false, ValidateUnknownRuleError
		}
	}
	return true, nil
}

func validateString(str string) (bool, error) {
	if emailRegexp.MatchString(str) ||
		ipRegexp.MatchString(str) ||
		usernameRegexp.MatchString(str) ||
		mobileRegexp.MatchString(str) ||
		telRegexp.MatchString(str) {
		return true, nil
	} else {
		return false, ValidateInvalidError
	}
}
