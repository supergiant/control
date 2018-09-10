package digitalocean

import (
	"testing"
	"strings"
)

func TestFingerPrint(t *testing.T) {
	expected := "ed:79:fd:40:e6:a1:05:64:ce:84:40:94:72:eb:9c:ee"
	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCpJTpadNp+c8MMQ/cyiWMjio5WxsFklDxG4RCuP6tgUEWqnANelNxT/lkIO5hUCfCS8a4wGPiOWIJpYMmmQRz7lysqm3hGGLVSv1H8m9XY//t/Xd+On7M/FZtr1AB/WV/11YBU8jW0TWk/pgPHjUUYnbbPAK1iilQS1ULx/Wen6EmjzVqD8XDLl82/cQgfT6UF1ZVQd+7qPmdeK4her+Otg/rTwIqjQI7DObhThpn7ZHehclTULw0jtAGw7/3Bek/DAuKSG3yQ+hMg+0xqO1t6zo12kYlRwpGTiCW2zLAuVw7PW7nz3SGvOTAjXAzKYcVdCn9rSs6UqufP4FV2BlbW3ZoQJY2KoEuDFbgmyhP8Z/+A6EXVkQBY/jHHsJGWIZS1QGpSAbYEGubb/lKryw0k1nr4X+bmFeymuOSWdipYOv/b4nXUrI+qIAZIza7heSM5BuRqkvVO/SSqyNbrypWHmL8x+EVb0WiSLQqFh/VZKiW0cgZ2gWL+qYyHuKlTPXCa+vO3SpPVFyIKV6WlblrSeCpwC6dj94RSkQejOojXvUJ1eT504dU8zyDYgE5nAgxeJecnM5+5Kowb/Zi5ByIjAmRE8e7ST4C9g73sue3t5foJ6IItJtlVgIoP5W3GLbRJ8p8T5SQY7fIVR6BiUmWU9BR2XdWVi2sH/x1IW9meoQ== glebstepanov1992@gmail.com"

	fg, err := fingerprint(publicKey)

	if err != nil {
		t.Error(err)
	}

	if !strings.EqualFold(fg, expected) {
		t.Errorf("Wrong fingerprint expected %s actual %s", expected, fg)
	}
}
