package testtelegram

type FakeTelegramApi struct {
}

func New() *FakeTelegramApi {
	return &FakeTelegramApi{
	}
}

func (f *FakeTelegramApi) Notify() chan string {
	return make(chan string)
}