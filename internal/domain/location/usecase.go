package location

type WebhookQueue interface {
	Enqueue(loc *Location)
}
