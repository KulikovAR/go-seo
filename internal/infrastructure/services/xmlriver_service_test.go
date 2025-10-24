package services

import (
	"encoding/xml"
	"testing"
)

func TestSearchResponseErrorHandling(t *testing.T) {
	// Тестовый XML с ошибкой
	xmlWithError := `<?xml version="1.0" encoding="UTF-8"?>
<yandexsearch version="1.0">
<request>
<query>вконтакте</query>
<page>1</page>
</request>
<response date="20251020T125529">
<error code="110">В данный момент сервис сильно перегружен. Попробуйте повторить запрос еще раз.</error>
</response>
</yandexsearch>`

	var searchResp SearchResponse
	err := xml.Unmarshal([]byte(xmlWithError), &searchResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Проверяем, что ошибка правильно парсится
	if searchResp.Response.Error == nil {
		t.Fatal("Expected error to be present in response")
	}

	if searchResp.Response.Error.Code != "110" {
		t.Errorf("Expected error code '110', got '%s'", searchResp.Response.Error.Code)
	}

	expectedMessage := "В данный момент сервис сильно перегружен. Попробуйте повторить запрос еще раз."
	if searchResp.Response.Error.Message != expectedMessage {
		t.Errorf("Expected error message '%s', got '%s'", expectedMessage, searchResp.Response.Error.Message)
	}
}

func TestSearchResponseSuccessHandling(t *testing.T) {
	// Тестовый XML без ошибки
	xmlWithoutError := `<?xml version="1.0" encoding="UTF-8"?>
<yandexsearch version="1.0">
<request>
<query>вконтакте</query>
<page>1</page>
</request>
<response date="20251020T125529">
<found>1000</found>
<results>
<grouping>
<group id="1" doccount="10">
<doc>
<url>https://vk.com</url>
<title>ВКонтакте</title>
<contenttype>organic</contenttype>
</doc>
</group>
</grouping>
</results>
</response>
</yandexsearch>`

	var searchResp SearchResponse
	err := xml.Unmarshal([]byte(xmlWithoutError), &searchResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Проверяем, что ошибки нет
	if searchResp.Response.Error != nil {
		t.Fatal("Expected no error in successful response")
	}

	// Проверяем, что данные правильно парсятся
	if searchResp.Response.Found != 1000 {
		t.Errorf("Expected found count 1000, got %d", searchResp.Response.Found)
	}

	if len(searchResp.Response.Results.Grouping.Groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(searchResp.Response.Results.Grouping.Groups))
	}
}
