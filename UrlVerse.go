package main

import (
        "bufio"
        "flag"
        "fmt"
        "io"
        "net/url"
        "os"
        "path"
        "regexp"
        "strings"
)

// Equivalences explained in README.md
// Feel free to add specific words here that can be normalized to clean up the results.
// The left side should be a unique string (like FOO) — keep it not too long but not too short.  This is used internally as replacement
// This is used for replacements. The words will be compiled into a regex, so don't mess with it, but you can use it to your advantage.
var equivalences = map[string][]string{
        //"TESLA": {"model-3", "model-y", ...},
        // langcodes are too small and need to be treated seperately
}

// Consider surrounding numbers with special characters or ensuring they're a certain length.
var numberregex = regexp.MustCompile("\\d+(\\.\\d+)?")
var profilepageregex = regexp.MustCompile("(?i)/(u|user|profile|author|member|referral)s?/[^/]+/?")
var titleregex = regexp.MustCompile("[A-Za-z0-9.]-[A-Za-z0-9.]-[A-Za-z0-9.-]+")
var langregex = buildlangregex()
var uuidregex = regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
var hashregex = regexp.MustCompile("[a-zA-Z0-9]{32,40,64,128}")
var equivalenceregexes = buildquivalences()

// change it based on your targets
var exts = []string{".css", ".png", ".jpg", ".jpeg", ".svg", ".gif", ".mp3", ".mp4", ".rss", ".ttf", ".woff", ".woff2", ".eot", ".pdf", ".m4v", ".ogv", ".webm"}
var paths = []string{"wp-content", "blog", "blogs", "product", "doc", "docs", "support", "about", "contact", "faq", "terms", "privacy", "help", "assets", "images", "videos", "css", "js", "scripts", "static", "uploads"}



var langcodes = []string{"af", "af-ZA", "ar", "ar-AE", "ar-BH", "ar-DZ", "ar-EG", "ar-IQ", "ar-JO", "ar-KW", "ar-LB", "ar-LY", "ar-MA", "ar-OM", "ar-QA", "ar-SA", "ar-SY", "ar-TN", "ar-YE", "az", "az-AZ", "az-AZ", "be", "be-BY", "bg", "bg-BG", "bs-BA", "ca", "ca-ES", "cs", "cs-CZ", "cy", "cy-GB", "da", "da-DK", "de", "de-AT", "de-CH", "de-DE", "de-LI", "de-LU", "dv", "dv-MV", "el", "el-GR", "en", "en-AU", "en-BZ", "en-CA", "en-CB", "en-GB", "en-IE", "en-JM", "en-NZ", "en-PH", "en-TT", "en-US", "en-ZA", "en-ZW", "eo", "es", "es-AR", "es-BO", "es-CL", "es-CO", "es-CR", "es-DO", "es-EC", "es-ES", "es-ES", "es-GT", "es-HN", "es-MX", "es-NI", "es-PA", "es-PE", "es-PR", "es-PY", "es-SV", "es-UY", "es-VE", "et", "et-EE", "eu", "eu-ES", "fa", "fa-IR", "fi", "fi-FI", "fo", "fo-FO", "fr", "fr-BE", "fr-CA", "fr-CH", "fr-FR", "fr-LU", "fr-MC", "gl", "gl-ES", "gu", "gu-IN", "he", "he-IL", "hi", "hi-IN", "hr", "hr-BA", "hr-HR", "hu", "hu-HU", "hy", "hy-AM", "id", "id-ID", "is", "is-IS", "it", "it-CH", "it-IT", "ja", "ja-JP", "ka", "ka-GE", "kk", "kk-KZ", "kn", "kn-IN", "ko", "ko-KR", "kok", "kok-IN", "ky", "ky-KG", "lt", "lt-LT", "lv", "lv-LV", "mi", "mi-NZ", "mk", "mk-MK", "mn", "mn-MN", "mr", "mr-IN", "ms", "ms-BN", "ms-MY", "mt", "mt-MT", "nb", "nb-NO", "nl", "nl-BE", "nl-NL", "nn-NO", "ns", "ns-ZA", "pa", "pa-IN", "pl", "pl-PL", "ps", "ps-AR", "pt", "pt-BR", "pt-PT", "qu", "qu-BO", "qu-EC", "qu-PE", "ro", "ro-RO", "ru", "ru-RU", "sa", "sa-IN", "se", "se-FI", "se-FI", "se-FI", "se-NO", "se-NO", "se-NO", "se-SE", "se-SE", "se-SE", "sk", "sk-SK", "sl", "sl-SI", "sq", "sq-AL", "sr-BA", "sr-BA", "sr-SP", "sr-SP", "sv", "sv-FI", "sv-SE", "sw", "sw-KE", "syr", "syr-SY", "ta", "ta-IN", "te", "te-IN", "th", "th-TH", "tl", "tl-PH", "tn", "tn-ZA", "tr", "tr-TR", "tt", "tt-RU", "ts", "uk", "uk-UA", "ur", "ur-PK", "uz", "uz-UZ", "uz-UZ", "vi", "vi-VN", "xh", "xh-ZA", "zh", "zh-CN", "zh-HK", "zh-MO", "zh-SG", "zh-TW", "zu", "zu-zA"}

var dotstar = regexp.MustCompile(".*")
var paramregexes = map[string]*regexp.Regexp{
        "utm_source":   dotstar,
        "utm_medium":   dotstar,
        "utm_campaign": dotstar,
        "utm_content":  dotstar,
        "utm_term":     dotstar,
        "redirect":     regexp.MustCompile("no"),
         // TODO: Consider adding version, v, cb, cache, etc.
}

func buildlangregex() *regexp.Regexp {
        // langcodes are currently only completely matched, might remove ^ and $ for the longer ones?
        reg := "(?i)^("
        for i, lang := range langcodes {
                reg += strings.Replace(lang, "-", "[-_]", 1) // match en-US & en_US
                if i < len(langcodes)-1 {
                        reg += "|"
                }
        }
        reg += ")$"
        return regexp.MustCompile(reg)
}

func buildquivalences() map[string]*regexp.Regexp {
        regexes := map[string]*regexp.Regexp{}
        for replacement, eqwords := range equivalences {
                regexes[replacement] = buildeqregex(eqwords)
        }
        return regexes
}

func buildeqregex(equivalentwords []string) *regexp.Regexp {
        reg := "("
        for i, word := range equivalentwords {
                reg += strings.Replace(word, "-", "[-_]", 1)
                if i < len(equivalentwords)-1 {
                        reg += "|"
                }
        }
        reg += ")"
        return regexp.MustCompile(reg)
}

func main() {
        printNormalized := flag.Bool("print-normalized", false, "print the normalized version of the urls (for debugging)")
        flag.Usage = func() {
                fmt.Printf("%s [OPTIONS] < urls.txt > less_urls.txt\n", os.Args[0])
                flag.PrintDefaults()
        }
        flag.Parse()
        runurlame(os.Stdin, os.Stdout, *printNormalized)
}

func runurlame(reader io.Reader, output io.ReadWriter, printNormalized bool) error {
        seen := map[string]bool{}
        stdin := bufio.NewScanner(reader)
        for stdin.Scan() {
                urlstr := stdin.Text()
                if u, err := url.Parse(urlstr); err == nil && len(urlstr) > 1 {
                        if lamefiletype(u) || profilepage(u) || lamedir(u) {
                              // Skip URLs that are definitely boring
                                continue
                        }
                        normalized := normalizeURL(urlstr)
                        if seen[normalized] {
                                continue
                        } else {
                                seen[normalized] = true
                                seen[urldecode(normalized)] = true //TODO check if this breaks things
                        }
                        if printNormalized {
                                fmt.Fprintf(output, "%s\n", normalized)
                        } else {
                                fmt.Fprintf(output, "%s\n", urlstr)
                        }
                }
        }
        return nil
}

func lamefiletype(u *url.URL) bool {
        filetype := strings.ToLower(path.Ext(u.Path))
        for _, ext := range exts {
                if filetype == ext {
                        return true
                }
        }
        return false
}

func lamedir(u *url.URL) bool {
        for i, part := range strings.Split(u.Path, "/") {
                lower := strings.ToLower(part)
                for _, lamepath := range paths {
                        if i > 2 {
                                //this is so we match /en-US/blog but not /api/v1/edit/blog
                                return false
                        }
                        if lower == lamepath {
                                return true
                        }
                }
        }
        return false
}

func profilepage(u *url.URL) bool {
        if profilepageregex.MatchString(u.Path) {
                return true
        }
        return false
}

func normalizeURL(urlstr string) string {
        //this func if there's an error, it just return the original URL
        if u, err := url.Parse(urlstr); err == nil {
                newvals := url.Values{}
                for key := range u.Query() {
                        if !lameparam(key, u.Query().Get(key)) {
                                // ignoring boring params, if we see /foo and /foo?utm_source=bar we only list the first
                                newvals.Set(normalizeItem(key), "!-P-!")
                        }
                }
                return newURL(u, normalizePath(u.Path), newvals)
        }
        return urlstr
}

func lameparam(key, val string) bool {
        if paramregexes[key] != nil {
                return paramregexes[key].MatchString(val)
        }
        return false
}

func normalizePath(path string) string {
        normalized := ""
        split := strings.Split(path, "/")
        for _, part := range split {
                if strings.TrimSpace(part) == "" {
                        continue
                }
                normalized += "/" + normalizeItem(part)
        }
        return normalized
}

func normalizeItem(item string) string {
        orig := item
        item = applyequivalences(item)
        item = hashregex.ReplaceAllString(item, "!-H-!")
        item = uuidregex.ReplaceAllString(item, "!-U-!")
        item = langregex.ReplaceAllString(item, "!-L-!")
        if len(item) > 10 && titleregex.MatchString(item) {
                return "!-T-!"
        }
        if orig == item {
                // only apply `numberregex` if hash / UUID wasn't found, might be too generic otherwise
                item = numberregex.ReplaceAllString(item, "!-N-!")
        }
        return item
}

func applyequivalences(item string) string {
        for replacement, regex := range equivalenceregexes {
                item = regex.ReplaceAllString(item, "!-"+replacement+"-!")
        }
        return item
}

func newURL(old *url.URL, path string, vals url.Values) string {
        return /*ignore scheme*/ cleanHostname(old) + path + "?" + vals.Encode() + "#" + old.Fragment
}

func cleanHostname(u *url.URL) string {
        if u.Port() == "80" || u.Port() == "443" {
                return u.Hostname()
        }
        return u.Host
}

func urldecode(str string) string {
        if decoded, err := url.QueryUnescape(str); err == nil {
                return decoded
        }
        return str
}
