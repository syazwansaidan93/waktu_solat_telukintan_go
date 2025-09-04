package main

import (
        "encoding/json"
        "fmt"
        "io"
        "net/http"
        "os"
        "time"
)

type PrayerTime struct {
        Date    string `json:"date"`
        Fajr    string `json:"fajr"`
        Dhuhr   string `json:"dhuhr"`
        Asr     string `json:"asr"`
        Maghrib string `json:"maghrib"`
        Isha    string `json:"isha"`
}

type PrayerAPIResponse struct {
        PrayerTime []PrayerTime `json:"prayerTime"`
}

type Prayer struct {
        Name          string `json:"name"`
        Time          string `json:"time"`
        IsHighlighted bool   `json:"isHighlighted"`
}

type CombinedData struct {
        CurrentDate string       `json:"currentDate"`
        CurrentTime string       `json:"currentTime"`
        PrayerTimes *PrayerTimes `json:"prayerTimes"`
}

type PrayerTimes struct {
        Date    string   `json:"date"`
        Prayers []Prayer `json:"prayers"`
}

func formatDisplayDate(date time.Time) string {
        monthsMalay := map[int]string{
                1:  "Januari",
                2:  "Februari",
                3:  "Mac",
                4:  "April",
                5:  "Mei",
                6:  "Jun",
                7:  "Julai",
                8:  "Ogos",
                9:  "September",
                10: "Oktober",
                11: "November",
                12: "Disember",
        }
        weekdaysMalay := map[time.Weekday]string{
                time.Sunday:    "Ahad",
                time.Monday:    "Isnin",
                time.Tuesday:   "Selasa",
                time.Wednesday: "Rabu",
                time.Thursday:  "Khamis",
                time.Friday:    "Jumaat",
                time.Saturday:  "Sabtu",
        }
        weekday := weekdaysMalay[date.Weekday()]
        month := monthsMalay[int(date.Month())]
        return fmt.Sprintf("%s, %d %s %d", weekday, date.Day(), month, date.Year())
}

func formatTime12HourNoAmpm(timeString24hr string) string {
        t, err := time.Parse("15:04:05", timeString24hr)
        if err != nil {
                return timeString24hr
        }
        return t.Format("3:04")
}

func fetchAndSaveMonthlyPrayerTimes(zoneCode string) error {
        prayerApiURL := fmt.Sprintf("https://www.e-solat.gov.my/index.php?r=esolatApi/takwimsolat&zone=%s&period=month", zoneCode)

        resp, err := http.Get(prayerApiURL)
        if err != nil {
                return fmt.Errorf("failed to fetch prayer times: %w", err)
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
        }

        body, err := io.ReadAll(resp.Body)
        if err != nil {
                return fmt.Errorf("failed to read response body: %w", err)
        }

        var apiResponse PrayerAPIResponse
        if err := json.Unmarshal(body, &apiResponse); err != nil {
                return fmt.Errorf("failed to decode API response: %w", err)
        }

        file, err := os.Create("prayer_times_monthly.json")
        if err != nil {
                return fmt.Errorf("failed to create file: %w", err)
        }
        defer file.Close()

        if err := json.NewEncoder(file).Encode(apiResponse.PrayerTime); err != nil {
                return fmt.Errorf("failed to write to file: %w", err)
        }

        return nil
}

func getTodaysPrayerTimes() (*PrayerTimes, error) {
        file, err := os.Open("prayer_times_monthly.json")
        if err != nil {
                return nil, fmt.Errorf("failed to open prayer times file: %w", err)
        }
        defer file.Close()

        var monthlyData []PrayerTime
        if err := json.NewDecoder(file).Decode(&monthlyData); err != nil {
                return nil, fmt.Errorf("failed to decode JSON from file: %w", err)
        }

        now := time.Now()
        var rawPrayerData *PrayerTime
        for _, dayData := range monthlyData {
                apiDate, err := time.Parse("02-Jan-2006", dayData.Date)
                if err == nil && apiDate.Year() == now.Year() && apiDate.Month() == now.Month() && apiDate.Day() == now.Day() {
                        rawPrayerData = &dayData
                        break
                }
        }

        if rawPrayerData == nil {
                return nil, fmt.Errorf("no prayer time data found for today")
        }

        prayerOrder := []string{"fajr", "dhuhr", "asr", "maghrib", "isha"}
        prayerNamesMalay := map[string]string{
                "fajr": "Subuh", "dhuhr": "Zuhur", "asr": "Asar", "maghrib": "Maghrib", "isha": "Isyak",
        }

        var allPrayerTimesWithDatetime []struct {
                Name          string
                FormattedTime string
                Datetime      time.Time
        }

        for _, key := range prayerOrder {
                var time24hr string
                switch key {
                case "fajr":
                        time24hr = rawPrayerData.Fajr
                case "dhuhr":
                        time24hr = rawPrayerData.Dhuhr
                case "asr":
                        time24hr = rawPrayerData.Asr
                case "maghrib":
                        time24hr = rawPrayerData.Maghrib
                case "isha":
                        time24hr = rawPrayerData.Isha
                }

                if time24hr == "" {
                        continue
                }

                parsedTime, err := time.Parse("15:04:05", time24hr)
                if err != nil {
                        continue
                }

                prayerDatetime := time.Date(now.Year(), now.Month(), now.Day(), parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), 0, now.Location())

                allPrayerTimesWithDatetime = append(allPrayerTimesWithDatetime, struct {
                        Name          string
                        FormattedTime string
                        Datetime      time.Time
                }{
                        Name:          prayerNamesMalay[key],
                        FormattedTime: formatTime12HourNoAmpm(time24hr),
                        Datetime:      prayerDatetime,
                })
        }

        var prayersOutput []Prayer
        nextPrayerIndex := -1
        for i, prayer := range allPrayerTimesWithDatetime {
                if prayer.Datetime.After(now) {
                        nextPrayerIndex = i
                        break
                }
        }

        if nextPrayerIndex != -1 {
                for i, p := range allPrayerTimesWithDatetime {
                        prayersOutput = append(prayersOutput, Prayer{
                                Name:          p.Name,
                                Time:          p.FormattedTime,
                                IsHighlighted: i == nextPrayerIndex,
                        })
                }
        } else if len(allPrayerTimesWithDatetime) > 0 {
                for i, p := range allPrayerTimesWithDatetime {
                        prayersOutput = append(prayersOutput, Prayer{
                                Name:          p.Name,
                                Time:          p.FormattedTime,
                                IsHighlighted: i == 0,
                        })
                }
        }

        return &PrayerTimes{
                Date:    formatDisplayDate(now),
                Prayers: prayersOutput,
        }, nil
}

func prayerTimesHandler(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Content-Type", "application/json")

        location, err := time.LoadLocation("Asia/Kuala_Lumpur")
        if err != nil {
                http.Error(w, `{"error":"Failed to load timezone"}`, http.StatusInternalServerError)
                return
        }
        now := time.Now().In(location)

        prayerTimes, err := getTodaysPrayerTimes()
        if err != nil {
                prayerTimes = nil
        }

        combinedData := CombinedData{
                CurrentDate: formatDisplayDate(now),
                CurrentTime: now.Format("3:04:05"),
                PrayerTimes: prayerTimes,
        }

        if err := json.NewEncoder(w).Encode(combinedData); err != nil {
                http.Error(w, `{"error":"Failed to encode JSON response"}`, http.StatusInternalServerError)
        }
}

func main() {
        go func() {
                zoneCode := "PRK05"
                if err := fetchAndSaveMonthlyPrayerTimes(zoneCode); err != nil {
                }

                ticker := time.NewTicker(30 * 24 * time.Hour)
                defer ticker.Stop()
                for range ticker.C {
                        if err := fetchAndSaveMonthlyPrayerTimes(zoneCode); err != nil {
                        }
                }
        }()

        http.HandleFunc("/api/prayer-times", prayerTimesHandler)

        if err := http.ListenAndServe(":502", nil); err != nil {
        }
}
