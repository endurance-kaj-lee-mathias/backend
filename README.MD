## Analyse: Encryptiestrategieën voor PII-gegevens

Bij het ontwerpen van een GDPR-conforme opslag van persoonsgegevens (PII) hebben we verschillende encryptiestrategieën onderzocht. Het doel was een oplossing te vinden die **veilig**, **onderhoudbaar** en **praktisch implementeerbaar** is binnen een server-side Go applicatie.

### 1. Doelstellingen
De gekozen encryptie-oplossing moet:
- persoonsgegevens beschermen tegen datalekken
- compatibel zijn met wachtwoordwijzigingen
- per gebruiker geïsoleerde data bieden
- zoekfunctionaliteit op bepaalde velden mogelijk maken
- voldoen aan gangbare security best practices

---

## 2. Onderzochte encryptie-opties

### 2.1 Eén globale encryptiesleutel
Bij deze aanpak wordt alle data met één centrale AES-sleutel geëncrypteerd.

**Voordelen**
- Eenvoudig te implementeren
- Weinig key management complexiteit

**Nadelen**
- Bij compromittering van de sleutel is **alle gebruikersdata blootgesteld**
- Geen isolatie per gebruiker
- Slechte schaalbaarheid vanuit security-oogpunt

**Conclusie**  
Afgewezen wegens te grote impact bij key leakage.

---

### 2.2 Encryptiesleutel afgeleid van gebruikersgegevens
In deze aanpak wordt de encryptiesleutel afgeleid van gebruikersinformatie, meestal het wachtwoord (via een KDF zoals PBKDF2 of Argon2).

**Voordelen**
- Geen aparte opslag van encryptiesleutels nodig
- Sleutel is uniek per gebruiker

**Nadelen**
- Wachtwoordwijzigingen maken bestaande data onleesbaar
- Encryptie en authenticatie raken sterk verweven
- Verhoogd risico op dataverlies
- Complexe migratie vereist bij sleutelwijzigingen

**Conclusie**  
Niet geschikt voor server-side encryptie van persistente data.

---

### 2.3 Client-side encryptie (password-based)
Data wordt op de client geëncrypteerd voordat deze naar de server wordt gestuurd.

**Voordelen**
- Server heeft geen toegang tot plaintext data
- Hoge privacygaranties

**Nadelen**
- Complexe UX (key recovery, password reset)
- Niet geschikt voor server-side verwerking of zoekfunctionaliteit
- Buiten scope van deze applicatie

**Conclusie**  
Afgewezen wegens complexiteit en scope.

---

## 3. Gekozen oplossing: Master key + per-user encryptiesleutel

### 3.1 Architectuur
Onze uiteindelijke oplossing maakt gebruik van twee sleuteltypes:

- **Master key**
    - Centrale sleutel
    - Wordt veilig opgeslagen buiten de database (secret manager / environment variable)
    - Wordt uitsluitend gebruikt om per-user keys te encrypten

- **Per-user encryptiesleutel**
    - Random gegenereerde 32-byte sleutel
    - Uniek per gebruiker
    - Gebruikt voor encryptie van alle PII-velden van die gebruiker
    - Versleuteld opgeslagen in de database

### 3.2 Encryptie-algoritme
- **AES-256-GCM**
    - Authenticated encryption
    - Bescherming tegen manipulatie
    - Native ondersteund in Go (`crypto/aes`, `crypto/cipher`)

---

## 4. Waarom per-user keys niet afgeleid worden van de gebruiker

### 4.1 Gebruikersdata is geen geheim
Gegevens zoals user ID, e-mail of username zijn voorspelbaar en vaak zichtbaar.  
Encryptiesleutels moeten daarentegen **onvoorspelbaar en cryptografisch sterk** zijn.

### 4.2 Onafhankelijkheid van wachtwoorden
Door encryptiesleutels los te koppelen van wachtwoorden:
- kunnen gebruikers hun wachtwoord wijzigen zonder dat data her-encryptie nodig is
- vermijden we dataverlies
- blijft de encryptielaag stabiel en onderhoudbaar

### 4.3 Beperking van schade (blast radius)
Als één per-user sleutel gecompromitteerd raakt:
- wordt enkel de data van die gebruiker blootgesteld
- andere gebruikers blijven beschermd

Dit volgt het principe van **least privilege**.

---

## 5. Zoekbaarheid van geëncrypteerde velden
Voor velden waarop gezocht moet worden (zoals e-mailadres):
- wordt een **SHA-256 hash** opgeslagen naast de ciphertext
- zoekopdrachten gebeuren via de hash
- decryptie gebeurt pas na het ophalen van het record

Dit combineert veiligheid met functionele zoekmogelijkheden.

---

## 6. Conclusie
De combinatie van een **master key** en **random per-user encryptiesleutels** biedt de beste balans tussen veiligheid, flexibiliteit en onderhoudbaarheid. Deze aanpak voorkomt dataverlies bij wachtwoordwijzigingen, beperkt de impact van sleutelcompromittering en volgt gangbare industry best practices voor GDPR-conforme opslag van persoonsgegevens.