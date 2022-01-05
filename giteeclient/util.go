package giteeclient

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	sdk "github.com/opensourceways/go-gitee/gitee"
)

func genrateRGBColor() string {
	v := rand.New(rand.NewSource(time.Now().Unix()))
	return fmt.Sprintf("%02x%02x%02x", v.Intn(255), v.Intn(255), v.Intn(255))
}

// GenResponseWithReference generates response with reference to the original comment.
func GenResponseWithReference(e *sdk.NoteEvent, reply string) string {
	format := `
@%s , %s

<details>

%s

</details>
`

	details := `
In response to [this](%s):

%s
`
	c := e.GetComment()

	return fmt.Sprintf(
		format, e.GetCommenter(), reply,
		fmt.Sprintf(details, c.GetHtmlUrl(), strings.ReplaceAll(">"+c.GetBody(), "\n", "\n>")),
	)
}
