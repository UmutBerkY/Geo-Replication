CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT,
    content_long TEXT,
    author TEXT NOT NULL,
    region TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO articles (title, summary, content_long, author, region)
VALUES
-- 1. Yazılım Mühendisliği
('Yazılım Mühendisliği', 
 'Planlama, analiz, tasarım, test ve bakım süreçlerinin bütünüdür.',
 'Yazılım mühendisliği, bir yazılımın fikir aşamasından kullanım ömrü sonuna kadar geçen sürecin bilimsel yöntemlerle yönetilmesidir. Gereksinimlerin doğru tanımlanması, modüler mimari tasarımı, sürüm kontrol sistemlerinin (Git) kullanımı, test otomasyonu ve bakım süreçleri bu disiplinin temel taşlarıdır. Agile ve DevOps metodolojileri, yazılım ekiplerinin esnek ve sürekli entegrasyon odaklı çalışmasını sağlar. Mühendislik yaklaşımı, kod kalitesini ve yeniden kullanılabilirliği artırır, hata oranını düşürür. Ayrıca proje yönetimi, kullanıcı deneyimi (UX) ve güvenlik testleri gibi çok yönlü bileşenleri içerir. Yazılım mühendisleri, yalnızca kod değil, sistemin sürekliliğini, ölçeklenebilirliğini ve sürdürülebilirliğini de tasarlar.',
 'Seed Bot', 'eu'),

-- 2. Yapay Zeka Uzman Sistemler
('Yapay Zeka Uzman Sistemler',
 'İnsanın karar verme yeteneğini taklit eden bilgi tabanlı sistemlerdir.',
 'Uzman sistemler, belirli bir uzmanlık alanındaki bilgiyi kural tabanlı biçimde temsil ederek insan uzmanların karar süreçlerini taklit eder. Temel bileşenleri bilgi tabanı ve çıkarım motorudur. Bilgi tabanı, “Eğer … o halde …” biçiminde tanımlanan kurallardan oluşur. Çıkarım motoru ise bu kuralları çalıştırarak sonuç üretir. Tıp, mühendislik, hukuk gibi alanlarda teşhis, tahmin ve planlama için kullanılır. Modern uzman sistemler, bulanık mantık ve makine öğrenmesiyle birleşerek dinamik hale gelmiştir. Böylece sistemler, yeni verilerle kendini güncelleyebilir. Kullanıcı arayüzleriyle desteklenmiş bu sistemler, karar destek süreçlerinde insan hatasını azaltır ve kurumsal bilgi birikimini dijital ortama taşır.',
 'Seed Bot', 'eu'),

-- 3. Bulut Bilişim
('Bulut Bilişim',
 'Veri ve uygulamaların internet tabanlı altyapılarda barındırılmasıdır.',
 'Bulut bilişim, kullanıcıların yerel donanım yatırımı yapmadan uzaktaki veri merkezlerinden işlem gücü ve depolama hizmeti almasını sağlar. Temel modelleri IaaS, PaaS ve SaaS’tır. Amazon Web Services, Google Cloud ve Microsoft Azure gibi sağlayıcılar ölçeklenebilir çözümler sunar. Sanallaştırma teknolojileri, kaynakları paylaştırarak maliyet verimliliği sağlar. Mikro servis ve konteyner yapıları (Docker, Kubernetes) uygulamaların modüler biçimde yönetilmesini sağlar. Bulut, felaket kurtarma, yük dengeleme, otomatik yedekleme ve yüksek erişilebilirlik gibi avantajlar sunar. Güvenlik açısından IAM (Identity and Access Management) politikaları, veri şifreleme ve denetim mekanizmaları kritik rol oynar.',
 'Seed Bot', 'eu'),

-- 4. Oyun Programlama
('Oyun Programlama',
 'Fizik, yapay zeka ve grafik motorlarıyla etkileşimli dijital dünyalar yaratma sanatıdır.',
 'Oyun programlama, sanatsal ve teknik disiplinlerin birleşimidir. Unity, Unreal Engine ve Godot gibi oyun motorları sayesinde fizik hesaplamaları, ışık simülasyonu, animasyon ve yapay zeka bileşenleri tek platformda buluşur. Oyun döngüsü, kare hızını (FPS) optimize ederek oyuncu deneyimini belirler. C++, C#, Lua gibi diller yüksek performans sağlar. Fizik motorları (Havok, PhysX) çarpışma ve yerçekimi hesaplarını yapar. Yapay zeka, karakterlerin davranışlarını gerçekçi hale getirir. Ses motorları, çevresel atmosferi destekler. Günümüzde VR/AR teknolojileriyle etkileşimli oyun deneyimi gelişmektedir. Oyun programlama sadece eğlence değil, aynı zamanda eğitim, savunma ve simülasyon alanlarında da kullanılmaktadır.',
 'Seed Bot', 'eu'),

-- 5. Kablosuz Ağlar
('Kablosuz Ağlar',
 'Veri iletiminin fiziksel kablo olmadan elektromanyetik dalgalarla yapılmasıdır.',
 'Kablosuz ağlar, cihazlar arasında veri aktarımını radyo frekansları üzerinden gerçekleştirir. Wi-Fi (IEEE 802.11), Bluetooth, Zigbee ve 5G gibi teknolojiler, farklı hız ve menzil gereksinimlerine yanıt verir. Kablosuz haberleşmede anten tasarımı, frekans bandı seçimi ve sinyal modülasyonu kritik öneme sahiptir. 5G teknolojisi, milisaniye seviyesinde gecikme ve gigabit hızlar sunarak otonom araçlar ve IoT cihazları için devrim yaratmıştır. Ağ güvenliği WPA3, VPN ve AES şifreleme ile sağlanır. Yeni nesil ağlar, AI destekli kanal tahsisi ve enerji tasarruflu veri yönlendirme algoritmalarıyla daha akıllı hale gelmektedir.',
 'Seed Bot', 'eu'),

-- 6. Veri Yapıları
('Veri Yapıları',
 'Verilerin bellekte düzenlenme biçimi ve erişim yöntemlerini tanımlar.',
 'Veri yapıları, yazılım performansını doğrudan etkileyen temel bileşenlerdir. Diziler, bağlı listeler, yığınlar, kuyruklar, ağaçlar ve grafikler farklı ihtiyaçlara göre kullanılır. Hash tablolar hızlı erişim sağlarken, ağaç yapıları sıralı veriler için idealdir. Doğru veri yapısının seçimi, algoritmaların karmaşıklığını azaltır. Big-O notasyonu, işlem maliyetlerini analiz eder. Bellek yönetimi, işaretçi yapıları ve dinamik tahsis süreçleri düşük seviyeli dillerde kritik öneme sahiptir. Gerçek dünya uygulamalarında sosyal ağ bağlantıları grafik yapılarıyla, arama motorları ise ters indeksleme algoritmalarıyla çalışır. Veri yapıları, bilgisayar biliminin DNA’sıdır.',
 'Seed Bot', 'eu'),

-- 7. Mikroişlemciler
('Mikroişlemciler',
 'Bilgisayarların beyni olarak işlem ve kontrol görevlerini yerine getirir.',
 'Mikroişlemciler, milyonlarca transistörü tek bir çipte birleştirerek tüm hesaplama ve kontrol görevlerini yürütür. CPU çekirdeği, ALU, kontrol birimi ve önbellek yapılarından oluşur. Komut seti mimarileri (RISC, CISC), performans ve güç tüketimini belirler. ARM mimarisi enerji verimliliğiyle mobil cihazlarda öne çıkarken, x86 mimarisi yüksek performanslı masaüstü sistemlerde tercih edilir. Mikroişlemciler, veri yolları aracılığıyla çevresel birimlerle iletişim kurar. Gömülü sistemlerde mikrodenetleyiciler (MCU) düşük güç tüketimiyle görev odaklı çalışır. Paralel işlem, pipeline ve çok çekirdekli mimariler sayesinde modern işlemciler saniyede milyarlarca işlemi gerçekleştirebilir.',
 'Seed Bot', 'eu'),

-- 8. Görüntü İşleme
('Görüntü İşleme',
 'Dijital görsellerin bilgisayar algoritmalarıyla analiz edilmesidir.',
 'Görüntü işleme, dijital bir görüntünün matematiksel temsiline dayalı olarak bilgi çıkarımını amaçlar. Gürültü azaltma, kontrast artırma, kenar tespiti ve segmentasyon gibi teknikler temel taşlardır. Bilgisayarla görü sistemleri (Computer Vision), endüstriyel kalite kontrol, tıbbi teşhis ve güvenlik alanlarında kullanılır. Derin öğrenme tabanlı CNN mimarileri, insan benzeri nesne tanıma başarısına ulaşmıştır. Renk uzayları (RGB, HSV, YCbCr), histogram analizi ve morfolojik işlemler görüntü kalitesini artırır. OpenCV, Pillow ve TensorFlow gibi kütüphaneler bu alandaki uygulamaların temelini oluşturur.',
 'Seed Bot', 'eu'),

-- 9. Kriptoloji
('Kriptoloji',
 'Bilginin güvenli şekilde şifrelenmesi ve çözümlenmesi bilimidir.',
 'Kriptoloji, matematiksel temelli bir bilimdir ve gizlilik, bütünlük, doğrulama gibi güvenlik kavramlarını destekler. Simetrik şifreleme (AES), asimetrik şifreleme (RSA), karma fonksiyonları (SHA-256) ve dijital imzalar bu bilimin uygulamalarıdır. HTTPS, blockchain, e-imza sistemleri kriptolojik protokollere dayanır. Kuantum bilgisayarların yükselişi, klasik kriptografi yöntemlerini tehdit etmektedir; bu nedenle post-kuantum algoritmalar geliştirilmektedir. Kriptoloji yalnızca teknik değil, stratejik bir alandır; ulusal güvenlik, bankacılık ve siber savunma gibi kritik sektörlerde veri gizliliğinin temelini oluşturur.',
 'Seed Bot', 'eu'),

-- 10. Ağ ve Bilgi Güvenliği
('Ağ ve Bilgi Güvenliği',
 'Veri iletiminde gizlilik, bütünlük ve erişilebilirliğin korunmasıdır.',
 'Ağ güvenliği, dijital verinin yetkisiz erişim ve değişikliğe karşı korunmasıdır. Temel prensipler gizlilik, bütünlük ve erişilebilirliktir (CIA triadı). Güvenlik duvarları, IDS/IPS sistemleri ve VPN tünelleme teknolojileri ağ savunmasının temelini oluşturur. TLS/SSL protokolleri veri aktarımını şifreler. Siber saldırılar (phishing, DDoS, ransomware) arttıkça güvenlik farkındalığı kritik hale gelmiştir. Kurumsal düzeyde ISO 27001 standartları, güvenlik politikaları için çerçeve sunar. Zero Trust modeli, kimlik doğrulama ve erişim kontrolünü her seviyede zorunlu kılar. Günümüzde yapay zeka destekli tehdit tespiti sistemleriyle proaktif koruma sağlanmaktadır.',
 'Seed Bot', 'eu'),

-- 11. İşletim Sistemleri
('İşletim Sistemleri',
 'Donanım ve kullanıcı arasında aracı görevi görür.',
 'İşletim sistemleri, donanım kaynaklarını yöneten ve uygulamalara hizmet sağlayan yazılımlardır. Çekirdek, süreç yönetimi, bellek tahsisi, dosya sistemi ve donanım sürücüleri temel bileşenlerdir. Linux, Windows ve macOS, farklı mimarilerle çalışır. Sanallaştırma ve konteyner teknolojileri (Docker, KVM) sistem kaynaklarının verimli paylaşımını sağlar. İşletim sistemleri, sistem çağrıları aracılığıyla donanımla güvenli iletişim kurar. Gerçek zamanlı sistemler (RTOS) milisaniyelik tepki süreleriyle otomotiv ve savunma sanayinde kritik rol oynar.',
 'Seed Bot', 'eu'),

-- 12. Veri Madenciliği
('Veri Madenciliği',
 'Büyük veri kümelerinden anlamlı bilgi çıkarma sürecidir.',
 'Veri madenciliği, veriler arasındaki örüntüleri, ilişkileri ve eğilimleri ortaya çıkaran bir analiz sürecidir. Kümelenme, sınıflandırma, regresyon ve birliktelik kuralları en yaygın yöntemlerdir. Büyük veri teknolojileri (Hadoop, Spark) yüksek hacimli verilerde paralel analiz sağlar. Sağlık, finans, pazarlama ve üretim alanlarında tahmine dayalı karar desteği sunar. Makine öğrenmesiyle birleşerek anomali tespiti ve öneri sistemleri gibi gelişmiş çözümler oluşturur. Veri madenciliği, ham veriyi anlamlı bilgiye dönüştürmenin bilimidir.',
 'Seed Bot', 'eu'),

-- 13. Makine Öğrenme
('Makine Öğrenme',
 'Verilerden örüntü çıkararak karar verme yeteneği kazandıran AI alt dalıdır.',
 'Makine öğrenme, verilerden istatistiksel çıkarımlar yaparak bilgisayarların örüntüleri tanımasını sağlar. Denetimli, denetimsiz ve pekiştirmeli öğrenme olarak üç ana kategoriye ayrılır. Algoritmalar arasında Decision Tree, Random Forest, SVM ve Neural Network modelleri bulunur. Model başarısı, veri kalitesi, öznitelik seçimi ve hiperparametre optimizasyonuna bağlıdır. Günümüzde ML, tahminleme, yüz tanıma, ses analizi ve doğal dil işleme gibi birçok alanda uygulanmaktadır. Büyük veriyle birleştiğinde, otomasyonun ve kişiselleştirmenin anahtarı haline gelmiştir.',
 'Seed Bot', 'eu'),

-- 14. Derin Öğrenme
('Derin Öğrenme',
 'Yapay sinir ağlarının çok katmanlı yapısıyla veriden anlam çıkarma sürecidir.',
 'Derin öğrenme, çok katmanlı yapay sinir ağlarını kullanarak veriden yüksek düzeyde soyutlama yapabilen bir yöntemdir. CNN, RNN, LSTM ve Transformer mimarileri farklı veri tipleri için özelleşmiştir. Görüntü işleme, konuşma tanıma ve doğal dil işleme alanlarında insan seviyesine yakın başarılar elde edilmiştir. GPU tabanlı paralel hesaplama ve büyük veri setleri, derin öğrenmenin başarısını artırmıştır. Ancak model karmaşıklığı, açıklanabilirlik (explainable AI) ve enerji tüketimi hâlâ araştırma konularıdır. Derin öğrenme, geleceğin özerk sistemlerinin bilişsel temelini oluşturur.',
 'Seed Bot', 'eu');
