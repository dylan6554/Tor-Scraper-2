<div align='center'>

<h1>Tor Scraper</h1>
<h4> <span>  </span> <a href="https://github.com/dylan6554/Tor-Scraper-2"> Uygulama </a> 
</div>
  
### :star2: **Proje Hakkında**
  
Tor Scraper, Tor ağındaki (.onion) siteleri güvenli bir şekilde tarayan, sayfaların temizlenmiş HTML verisini çıkaran ve görsel kanıt için ekran görüntüsü alan bir otomasyon aracıdır. Elde edilen veriler bir web panel üzerinden kolayca izlenebilir.

### :dart: Özellikler
- Temizlenmiş HTML
- Ekran Görüntüsü
- Web Panel
## :toolbox: Teknolojiler
Bu proje aşağıdaki temel teknolojiler kullanılarak geliştirilmiştir:
- **Go (Golang)** 
- **PostgreSQL** 
- **Tor Network**
- **Docker & Docker Compose**

## :brain: Başlık Üretim Mantığı
Uygulama, taranan sayfanın içeriğini analiz ederek anlamlı başlıklar üretir:
Başlık (Title): Sayfa içeriği analiz edilerek en üst seviye başlık olan `<h1>` etiketi (`document.querySelector("h1")`) ile içindeki metin çekilir.



### :bangbang: Önkoşullar

- Bilgisayarınıza Docker`ı kurun.<a href="https://www.docker.com/products/docker-desktop/"> Buradan</a>


### :gear: Kurulum

Programı indirdikten sonra terminalde aynı klasöre gelip komutu çalıştırın.
```bash
docker compose up --build
```

Web panele giriş için user ve pass  `admin`  dir.
