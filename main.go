package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

var languages = []language.Tag{
	language.Afrikaans,
	language.Amharic,
	language.Arabic,
	language.ModernStandardArabic,
	language.Azerbaijani,
	language.Bulgarian,
	language.Bengali,
	language.Catalan,
	language.Czech,
	language.Danish,
	language.German,
	language.Greek,
	language.Spanish,
	language.EuropeanSpanish,
	language.LatinAmericanSpanish,
	language.Estonian,
	language.Persian,
	language.Finnish,
	language.Filipino,
	language.French,
	language.CanadianFrench,
	language.Gujarati,
	language.Hebrew,
	language.Hindi,
	language.Croatian,
	language.Hungarian,
	language.Armenian,
	language.Indonesian,
	language.Icelandic,
	language.Italian,
	language.Japanese,
	language.Georgian,
	language.Kazakh,
	language.Khmer,
	language.Kannada,
	language.Korean,
	language.Kirghiz,
	language.Lao,
	language.Lithuanian,
	language.Latvian,
	language.Macedonian,
	language.Malayalam,
	language.Mongolian,
	language.Marathi,
	language.Malay,
	language.Burmese,
	language.Nepali,
	language.Dutch,
	language.Norwegian,
	language.Punjabi,
	language.Polish,
	language.Portuguese,
	language.BrazilianPortuguese,
	language.EuropeanPortuguese,
	language.Romanian,
	language.Russian,
	language.Sinhala,
	language.Slovak,
	language.Slovenian,
	language.Albanian,
	language.Serbian,
	language.SerbianLatin,
	language.Swedish,
	language.Swahili,
	language.Tamil,
	language.Telugu,
	language.Thai,
	language.Turkish,
	language.Ukrainian,
	language.Urdu,
	language.Uzbek,
	language.Vietnamese,
	language.Chinese,
	language.SimplifiedChinese,
	language.TraditionalChinese,
	language.Zulu,
}

func main() {
	var (
		apiKey = flag.String("api-key", "", "A Google Translate API key")
		hops   = flag.Int("hops", 5, "the number of languages to pass the text through")
	)
	flag.Parse()

	text := flag.Arg(0)
	if text == "" {
		fmt.Fprintln(os.Stderr, "usage: translate-game -api-key <translate key> -hops 5 <phrase>")
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())
	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(*apiKey))
	if err != nil {
		log.Fatal(err)
	}

	if *hops > len(languages) {
		*hops = len(languages)
	}
	if *hops <= 0 {
		*hops = 1
	}

	prevLang := language.English

	var cycle []language.Tag
	for _, i := range rand.Perm(len(languages))[:*hops] {
		cycle = append(cycle, languages[i])
	}
	cycle = append(cycle, language.English)

	for _, lang := range cycle {
		text = strings.Replace(text, "\n", "<br />", -1)

		resp, err := client.Translate(ctx, []string{text}, lang, &translate.Options{
			Source: prevLang,
			Format: translate.HTML,
		})
		if e, ok := err.(*googleapi.Error); ok {
			if e.Code == 400 {
				continue
			}
		}
		if err != nil {
			log.Fatal(err)
		}

		translation := resp[0]
		text = translation.Text
		prevLang = lang

		text = strings.Replace(text, "<br />", "\n", -1)
		fmt.Println(text)
	}
}
