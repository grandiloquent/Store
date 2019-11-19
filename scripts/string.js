const gen = function () {
    var str = [];
    for (let i = 0; i < 26; i++) {
        str.push(`deauthPacket[${i}]`)
    }

    let a = 'Serial.printf("' + "%02x".repeat(26) + '\\n",' + str.join(",") + ');';
    console.log(a);
}

gen();