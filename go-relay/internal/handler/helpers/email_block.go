package helpers

import (
	"net/mail"
	"os"
	"strings"
)

// ValidateEmailForRegistration checks email format and blocks disposable domains.
// Returns an error message to send to the client, or "" if the email is acceptable.
func ValidateEmailForRegistration(email string) string {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "Invalid email address."
	}
	parts := strings.SplitN(addr.Address, "@", 2)
	if len(parts) != 2 || !strings.Contains(parts[1], ".") {
		return "Invalid email address."
	}
	domain := strings.ToLower(parts[1])
	if isBlockedDomain(domain) {
		return "This email domain is not allowed."
	}
	return ""
}

func isBlockedDomain(domain string) bool {
	if disposableEmailDomains[domain] {
		return true
	}
	extra := os.Getenv("BLOCKED_EMAIL_DOMAINS")
	if extra == "" {
		return false
	}
	for _, d := range strings.Split(extra, ",") {
		if strings.EqualFold(strings.TrimSpace(d), domain) {
			return true
		}
	}
	return false
}

// disposableEmailDomains is a curated list of common throwaway email providers.
// Operators can extend this via the BLOCKED_EMAIL_DOMAINS env var (comma-separated).
var disposableEmailDomains = map[string]bool{
	"mailinator.com":       true,
	"guerrillamail.com":    true,
	"guerrillamail.net":    true,
	"guerrillamail.org":    true,
	"guerrillamail.biz":    true,
	"guerrillamail.de":     true,
	"guerrillamail.info":   true,
	"guerrillamailblock.com": true,
	"grr.la":               true,
	"sharklasers.com":      true,
	"spam4.me":             true,
	"tempmail.com":         true,
	"temp-mail.org":        true,
	"10minutemail.com":     true,
	"10minutemail.net":     true,
	"throwam.com":          true,
	"throwaway.email":      true,
	"trashmail.com":        true,
	"trashmail.me":         true,
	"trashmail.net":        true,
	"trashmail.at":         true,
	"trashmail.io":         true,
	"dispostable.com":      true,
	"maildrop.cc":          true,
	"mailnesia.com":        true,
	"mailnull.com":         true,
	"discard.email":        true,
	"mailsac.com":          true,
	"yopmail.com":          true,
	"yopmail.fr":           true,
	"cool.fr.nf":           true,
	"jetable.fr.nf":        true,
	"jetable.com":          true,
	"jetable.net":          true,
	"jetable.org":          true,
	"nospam.ze.tc":         true,
	"nomail.xl.cx":         true,
	"nomail.pw":            true,
	"nospamfor.us":         true,
	"mega.zik.dj":          true,
	"speed.1s.fr":          true,
	"courriel.fr.nf":       true,
	"moncourrier.fr.nf":    true,
	"monemail.fr.nf":       true,
	"monmail.fr.nf":        true,
	"spamgourmet.com":      true,
	"spamgourmet.net":      true,
	"spamgourmet.org":      true,
	"spamspot.com":         true,
	"spamevader.com":       true,
	"spamcorner.com":       true,
	"spamcero.com":         true,
	"spamcon.org":          true,
	"spamfree24.org":       true,
	"spamgob.com":          true,
	"spaml.com":            true,
	"spamoff.de":           true,
	"spamthis.co.uk":       true,
	"spamtrap.ro":          true,
	"spamboy.com":          true,
	"spam.la":              true,
	"spamfighter.cf":       true,
	"spamfighter.ga":       true,
	"spamfighter.gq":       true,
	"spamfighter.ml":       true,
	"spamfighter.tk":       true,
	"tradermail.info":      true,
	"objectmail.com":       true,
	"rklips.com":           true,
	"getonemail.net":       true,
	"mailfreeonline.com":   true,
	"mailguard.me":         true,
	"mail-temporaire.fr":   true,
	"tempr.email":          true,
	"tempe.email":          true,
	"tmail.com":            true,
	"tmail.io":             true,
	"mohmal.com":           true,
	"mailtemp.info":        true,
	"burnermail.io":        true,
	"getairmail.com":       true,
	"fakeinbox.com":        true,
	"inboxbear.com":        true,
	"tempinbox.com":        true,
	"emailondeck.com":      true,
	"filzmail.com":         true,
	"einrot.com":           true,
	"tempomail.fr":         true,
	"owlpic.com":           true,
	"inemail.com":          true,
	"inoutmail.de":         true,
	"inoutmail.eu":         true,
	"inoutmail.info":       true,
	"inoutmail.net":        true,
	"lroid.com":            true,
	"mailboxy.fun":         true,
	"mailscrap.com":        true,
	"mailslapping.com":     true,
	"mailslite.com":        true,
	"mailsiphon.com":       true,
	"sogetthis.com":        true,
	"soodonims.com":        true,
	"superrito.com":        true,
	"tempalias.com":        true,
	"tempail.com":          true,
	"tempem.com":           true,
	"tempsky.com":          true,
	"wibblesmith.com":      true,
}
