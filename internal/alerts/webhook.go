package alerts

type Platform string

const (
	PlatformSlack   Platform = "SLACK"
	PlatformDiscord Platform = "DISCORD"
)

// WebhookSender dispatches alerts to a single slack or discord webhook URL.
type WebhookSender struct {
	url      string
	platform Platform
	client   *http.Client
}

// Slack payload schema
type slackPayload struct {
	Text        string            `json:"text"`
	Attachments []slackAttachment `json:"attachments"`
}

type slackAttachment struct {
	Colour string       `json:"color"`
	Title  string       `json:"title"`
	Text   string       `json:"text"`
	Fields []slackField `json:"fields"`
	Footer string       `json:"footer"`
	Ts     int64        `json:"ts"`
}

type slackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Discord payload schema
type discordPayload struct {
	Username string         `json:"username"`
	Embeds   []discordEmbed `json:"embeds"`
}

type discordEmbed struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Color       int            `json:"color"`
	Fields      []discordField `json:"fields"`
	Footer      discordFooter  `json:"footer"`
	Timestamp   time.Time      `json:"timestamp"`
}

type discordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}
type discordFooter struct {
	Text string `json:"text"`
}

func NewWebhookSender(url string, platform Platform) *WebhookSender {
	return &WebhookSender{
		url:      url,
		platform: platform,
		client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// Send formats and POSTs the alert to the configured webhook URL.
func (s *WebhookSender) Send(alert *models.Alert) error {
	payload, err := s.buildPayload(alert)
	if err != nil {
		return fmt.Errorf("failed to build payload: %w", err)
	}

	resp, err := s.client.Post(s.url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("http post to  %s failed: %w", s.platform, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http post to %s failed: %s", s.platform, resp.Status)
	}

	return nil
}

func (s *WebhookSender) buildPayload(alert *models.Alert) ([]byte, error) {
	switch s.platform {
	case PlatformSlack:
		return json.Marshal(w.slackPayload(alert))
	case PlatformDiscord:
		return json.Marshal(w.discordPayload(alert))
	default:
		return nil, fmt.Errorf("unsupported platform: %s", s.platform)
	}
}
func (s *webhookSender) slackPayload(alert *models.Alert) slackPayload {
	return slackPayload{
		Text: fmt.Sprintf(":rotating_light: Supaspy Alert %s", alert.Severity),
		Attachments: []slackAttachment{{
			Colour: slackColour(alert.Severity),
			Title: alert.Title,
			Text: alert.Message,
			Fields: []slackField{
				{
				Title: "Query",
				Value : alert.Event.ID,
			    Short: true},
				{
				Title: "Duration",
				}
			},
			Footer: "Supaspy Observability",
			Ts: alert.Timestamp.Unix(),
		}},
		}
	}
