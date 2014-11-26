package emoji

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	TwemojiHTMLTemplate  = `<img src="%s" width="%d" height="%d" >`
	TwemojiXHTMLTemplate = `<img src="%s" width="%d" height="%d" />`
)

var (
	hE10 = regexp.MustCompile(`&#[0-9]+;`)
	hE16 = regexp.MustCompile(`&#x[A-Fa-f0-9]+;`)
	nR10 = regexp.MustCompile(`[0-9]+`)
	nR16 = regexp.MustCompile(`[A-Fa-f0-9]+`)
)

func code2entities(src []rune) string {
	ret := make([]string, 0)

	for _, char := range src {
		ret = append(ret, fmt.Sprintf(`&#x%X;`, char))
	}

	return strings.Join(ret, ``)
}

func hexlizeEntities(src string) string {
	return hE10.ReplaceAllStringFunc(src, func(str string) string {
		num, err := strconv.ParseInt(nR10.FindString(str), 10, 32)
		if err != nil {
			return str
		}

		char := rune(num)
		if utf8.ValidRune(char) {
			return code2entities([]rune{char})
		}

		return str
	})
}

func EmojiTagToHTMLEntities(src string) string {
	for name, chars := range name2codes {
		tag := strings.Join([]string{`:`, name, `:`}, ``)

		src = strings.Replace(src, tag, code2entities(chars), -1)
	}

	return src
}

func EmojiTagToUnicode(src string) string {
	for name, chars := range name2codes {
		str := string(chars)
		tag := strings.Join([]string{`:`, name, `:`}, ``)

		src = strings.Replace(src, tag, str, -1)
	}

	return src
}

func EmojiTagToTwemoji(src string, size int, isXHTML bool) string {
	for name, chars := range name2codes {
		str := string(chars)
		if img, ok := str2img[str]; ok {
			var tpl string
			if isXHTML {
				tpl = TwemojiXHTMLTemplate
			} else {
				tpl = TwemojiHTMLTemplate
			}

			imgTag := fmt.Sprintf(tpl, img, size, size)
			tagStr := strings.Join([]string{`:`, name, `:`}, ``)

			src = strings.Replace(src, tagStr, imgTag, -1)
		}
	}

	return src
}

func UnicodeToHTMLEntities(src string) string {
	for _, chars := range name2codes {
		str := string(chars)
		entities := code2entities(chars)

		src = strings.Replace(src, str, entities, -1)

		chars2 := make([]rune, 0)
		for _, char := range chars {
			if char == '\uFE0F' {
				continue
			}

			chars2 = append(chars2, char)
		}

		str = string(chars2)
		entities = code2entities(chars2)

		src = strings.Replace(src, str, entities, -1)
	}

	return src
}

func UnicodeToTwemoji(src string, size int, isXHTML bool) string {
	for _, chars := range name2codes {
		keyStr := string(chars)
		img, ok := str2img[keyStr]

		if ok {
			str := string(chars)
			var tpl string
			if isXHTML {
				tpl = TwemojiXHTMLTemplate
			} else {
				tpl = TwemojiHTMLTemplate
			}

			tag := fmt.Sprintf(tpl, img, size, size)

			src = strings.Replace(src, str, tag, -1)

			chars2 := make([]rune, 0)
			for _, char := range chars {
				if char == '\uFE0F' {
					continue
				}

				chars2 = append(chars2, char)
			}

			str = string(chars2)
			src = strings.Replace(src, str, tag, -1)

		}
	}

	return src
}

func HTMLEntitiesToUnicode(src string) string {
	src = hexlizeEntities(src)
	for _, chars := range name2codes {
		str := string(chars)
		entities := code2entities(chars)

		src = strings.Replace(src, entities, str, -1)

		chars2 := make([]rune, 0)
		for _, char := range chars {
			if char == '\uFE0F' {
				continue
			}

			chars2 = append(chars2, char)
		}

		str = string(chars2)
		entities = code2entities(chars2)

		src = strings.Replace(src, entities, str, -1)
	}

	return src
}
