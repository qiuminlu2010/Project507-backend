package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

type App struct {
	JwtSecret       string
	PageSize        int
	RuntimeRootPath string
	ImagePrefixUrl  string
	VideoMaxSize    int

	ThumbMaxHeight int
	ThumbMaxWidth  int
	ImageMaxSize   int
	ImageAllowExts []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string

	AdminId uint
}

var AppSetting = &App{}

type Minio struct {
	Host              string
	EndPoint          string
	AccessKeyID       string
	SecretAccessKey   string
	ImageBucketName   string
	VideoBucketName   string
	PreviewBucketName string
	TempBucketName    string
	AvatarBucketName  string
	// ThumbnailBucketName    string
	// AvatarSavePath        string
	// ImageSavePath         string
	// ImageTempSavePath     string
	// ThumbSavePath         string
	// VideoTempSavePath     string
	// VideoSavePath         string
	// VideoPreviewSavePath  string
	// VideoCompressSavePath string
}

var MinioSetting = &Minio{}

type Server struct {
	RunMode      string
	HttpPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var ServerSetting = &Server{}

type Database struct {
	Type       string
	User       string
	Password   string
	HostMaster string
	HostSlave1 string
	HostSlave2 string
	// Host        string
	Name        string
	TablePrefix string
}

var DatabaseSetting = &Database{}

type Redis struct {
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	Host        string
	Password    string
}

var RedisSetting = &Redis{}

type Nsq struct {
	NsqLookup string
	Nsqd      string
}

var NsqSetting = &Nsq{}

func Setup() {
	Cfg, err := ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}

	err = Cfg.Section("app").MapTo(AppSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}
	err = Cfg.Section("minio").MapTo(MinioSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo AppSetting err: %v", err)
	}

	AppSetting.ImageMaxSize = AppSetting.ImageMaxSize * 1024 * 1024
	AppSetting.VideoMaxSize = AppSetting.VideoMaxSize * 1024 * 1024
	err = Cfg.Section("server").MapTo(ServerSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo ServerSetting err: %v", err)
	}

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second

	err = Cfg.Section("database").MapTo(DatabaseSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo DatabaseSetting err: %v", err)
	}
	// log.Println("DatabaseSetting", DatabaseSetting)
	// log.Println("DatabaseSetting.Host", DatabaseSetting.Host)

	err = Cfg.Section("redis").MapTo(RedisSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}

	err = Cfg.Section("nsq").MapTo(NsqSetting)
	if err != nil {
		log.Fatalf("Cfg.MapTo NsqSetting err: %v", err)
	}
}
