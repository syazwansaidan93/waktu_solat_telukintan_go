<?php
header("Access-Control-Allow-Origin: *");
header("Content-Type: application/json");

$zoneCode = 'PRK05';

date_default_timezone_set('Asia/Kuala_Lumpur');

function formatDisplayDate($dateObj) {
    $monthsMalay = [
        1 => 'Januari', 2 => 'Februari', 3 => 'Mac', 4 => 'April', 5 => 'Mei', 6 => 'Jun',
        7 => 'Julai', 8 => 'Ogos', 9 => 'September', 10 => 'Oktober', 11 => 'November', 12 => 'Disember'
    ];
    $weekdaysMalay = [
        0 => 'Ahad', 1 => 'Isnin', 2 => 'Selasa', 3 => 'Rabu', 4 => 'Khamis', 5 => 'Jumaat', 6 => 'Sabtu'
    ];

    $weekday = $weekdaysMalay[(int)$dateObj->format('w')];
    $month = $monthsMalay[(int)$dateObj->format('n')];
    return "{$weekday}, {$dateObj->format('j')} {$month} {$dateObj->format('Y')}";
}

function formatTime12HourNoAmpm($timeString24hr) {
    if (empty($timeString24hr)) {
        return '';
    }
    $timeObj = DateTime::createFromFormat('H:i:s', $timeString24hr); 
    if ($timeObj === false) {
        return $timeString24hr;
    }
    return $timeObj->format('g:i'); // Changed from 'h:i' to 'g:i' to remove leading zero
}

$currentFormattedDate = formatDisplayDate(new DateTime());
$currentFormattedTime = (new DateTime())->format('g:i:s'); // Changed from 'h:i:s' to 'g:i:s'

$prayerTimesData = null;
$allPrayerTimesForHighlight = [];

$prayerApiUrl = "https://www.e-solat.gov.my/index.php?r=esolatApi/takwimsolat&zone={$zoneCode}&period=today";

$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, $prayerApiUrl);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
curl_setopt($ch, CURLOPT_TIMEOUT, 30);
$prayerResponse = curl_exec($ch);
$httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);
curl_close($ch);

if ($prayerResponse !== false && $httpCode === 200) {
    $prayerJson = json_decode($prayerResponse, true);
    
    if ($prayerJson && isset($prayerJson['prayerTime'][0])) {
        $rawPrayerData = $prayerJson['prayerTime'][0];
        
        $prayerOrder = ['fajr', 'dhuhr', 'asr', 'maghrib', 'isha'];
        $prayerNamesMalay = [
            'fajr' => 'Subuh', 'dhuhr' => 'Zuhur', 'asr' => 'Asar',
            'maghrib' => 'Maghrib', 'isha' => 'Isyak'
        ];

        foreach ($prayerOrder as $key) {
            if (isset($rawPrayerData[$key])) {
                $time24hr = $rawPrayerData[$key];
                $formatted12hrNoAmpm = formatTime12HourNoAmpm($time24hr);
                
                $prayerDt = DateTime::createFromFormat('H:i:s', $time24hr); 
                if ($prayerDt === false) {
                    continue;
                }
                $prayerDt->setDate((new DateTime())->format('Y'), (new DateTime())->format('m'), (new DateTime())->format('d'));
                
                $allPrayerTimesForHighlight[] = [
                    'name' => $prayerNamesMalay[$key],
                    'time' => $formatted12hrNoAmpm,
                    'datetime_obj' => $prayerDt,
                    'isHighlighted' => false
                ];
            }
        }
        
        $nextPrayerIndex = -1;
        $minDiff = new DateInterval('P365D');
        $currentTime = new DateTime();

        foreach ($allPrayerTimesForHighlight as $index => $prayer) {
            $diff = $currentTime->diff($prayer['datetime_obj']);
            if ($prayer['datetime_obj'] > $currentTime) {
                $diffInSeconds = $diff->days * 86400 + $diff->h * 3600 + $diff->i * 60 + $diff->s;
                $minDiffInSeconds = $minDiff->days * 86400 + $minDiff->h * 3600 + $minDiff->i * 60 + $minDiff->s;

                if ($diffInSeconds < $minDiffInSeconds) {
                    $minDiff = $diff;
                    $nextPrayerIndex = $index;
                }
            }
        }
        
        if ($nextPrayerIndex === -1 && !empty($allPrayerTimesForHighlight)) {
            $nextPrayerIndex = 0;
        }

        if ($nextPrayerIndex !== -1) {
            $allPrayerTimesForHighlight[$nextPrayerIndex]['isHighlighted'] = true;
        }

        $prayersOutput = [];
        foreach ($allPrayerTimesForHighlight as $prayer) {
            $prayersOutput[] = [
                'name' => $prayer['name'],
                'time' => $prayer['time'],
                'isHighlighted' => $prayer['isHighlighted']
            ];
        }

        $prayerTimesData = [
            'date' => $currentFormattedDate,
            'prayers' => $prayersOutput
        ];
    }
}

$combinedData = [
    "currentDate" => $currentFormattedDate,
    "currentTime" => $currentFormattedTime,
    "prayerTimes" => $prayerTimesData,
];

echo json_encode($combinedData);
?>
