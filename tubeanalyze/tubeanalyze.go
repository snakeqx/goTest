package tubeanalyze

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//Global variables
const TimeFormat ="2006-01-02 15:04:05.000000"

type TubeHistory struct {
	TubeHistoryVersionType       string
	SoftwareVersionPrevious      string
	InstallationDateOfPreviousSW string
	SoftwareVersion              string
	SwInstallationDate           string
	TubeSystem                   string
	CustomerCityDistrict         string
	DateOfInstallation           string
	DateOfDeinstallation         string
	SystemSerialNumber           string
	TubeSerialNumber             string
	TubeRevision                 string
	TubeType                     string
	NumberOfScans                string
	KwsSinceTubeInst             string
	TubeScanSeconds              string
	SystemScanSeconds            string
	Scans                        *ScanHistory
}

type ScanHistory struct {
	TubeScanCounter          []int       //0
	Kind                     []string    //1
	DateTime                 []time.Time //2
	RotTime                  []float64   //3
	ScanTimeNom              []float64   //4
	ScanTimeAct              []float64   //5
	VoltageNom               []float64   //6
	VoltabeAct               []float64   //7
	CurrentUi                []int       //8
	CurrentNom               []int       //9
	CurrentMin               []int       //10
	CurrentMax               []int       //11
	CurrentMean              []int       //12
	CurrentCtrlBeg           []int       //13
	CurrentBeg               []int       //14
	CurrentEnd               []int       //15
	Focus                    []string    //16
	FreqStatorAct            []int       //17
	FreqAnodeAct             []int       //18
	FilamentCurrentNom       []int       //19
	FilamentCurrentBeg       []int       //20
	FilamentCurrentCtrlBeg   []int       //21
	FilamentCurrentEnd       []int       //22
	FilamentCurrentPushCalc  []int       //23
	DoseNom                  []int       //24
	DoseMin                  []int       //25
	DoseMax                  []int       //26
	DoseEnd                  []int       //27
	TubeCoolingTempInletBeg  []int       //28
	TubeCoolingTempInletEnd  []int       //29
	TubeCoolingTempOutletBeg []int       //30
	TubeCoolingTempOutletEnd []int       //31
	CalculatedTempAnode      []int       //32
	GantryTemperature        []int       //33
	StatorCurrentEndAct      []int       //34
	ArcingsSum               []int       //35
	ArcingsPos               []int       //36
	ArcingsNeg               []int       //37
	HvDrops                  []int       //38
	XcDrops                  []int       //39
	StartAngle               []int       //40
	ReadingsTotal            []int       //41
	ReadingsDefScan          []int       //42
	ReadinsLastDef           []int       //43
	Region                   []string    //44
	DOMType                  []string    //45
	ScanMode                 []string    //46
	CancelReason             []string    //47
	AbortController          []string    //48
	ECO                      []string    //49
	DiffToTUC                []string    //50
	TubeSerial               []int       //51
	SystemSerial             []int       //52
}

func NewTubeHistory(filename string) (*TubeHistory, error) {
	tube := new(TubeHistory)
	scans := new(ScanHistory)
	tube.Scans = scans

	err := tube.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return tube, nil

}

func (t *TubeHistory) ReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	countHeader := 0
	for scanner.Scan() {
		switch strings.HasPrefix(scanner.Text(), "#") {
		case true:
			err = t.handleHeader(scanner.Text(), countHeader)
			if err!=nil{
				countHeader = 0
				return err
			}
			countHeader += 1
		case false:
			err = t.handleData(scanner.Text())
			if err!=nil{
				return err
			}
		}
	}
	countHeader = 0

	return nil
}

func (t *TubeHistory) handleHeader(header string, headerIndex int) error {
	info := strings.SplitN(header, ":", 2)
	switch len(info) {
	case 1:
		info = strings.Split(info[0], "\t")
		if len(info) != 53 {
			return errors.New("the data format is not correct")
		}
	case 2:
		for i, s := range info {
			info[i] = strings.TrimSpace(s)
			// if i is even which means it is the content
			if i & 0x1 == 1 {
				ref := reflect.ValueOf(t).Elem()
				ref.Field(headerIndex).SetString(s)
			}
		}

	default:
		return errors.New("cannot parse the file, different format")
	}

	return nil
}

func (t *TubeHistory) handleData(data string) error {
	scans := strings.Split(data, "\t")
	ref := reflect.ValueOf(t.Scans).Elem()
	for index, scan := range scans {
		switch index {
		// case float
		case 3, 4, 5, 6, 7:
			f, err := strconv.ParseFloat(scan, 64)
			if err != nil {
				return err
			}
			value := reflect.ValueOf(f)
			refIndex := ref.Field(index)
			refIndex.Set(reflect.Append(refIndex, value))
		// case string
		case 1, 16, 44, 45, 46, 47, 48, 49, 50:
			value := reflect.ValueOf(scan)
			refIndex := ref.Field(index)
			refIndex.Set(reflect.Append(refIndex, value))
		// case time
		case 2:
			t, err := time.Parse(TimeFormat, scan)
			if err!=nil{
				return err
			}
			value := reflect.ValueOf(t)
			refIndex := ref.Field(index)
			refIndex.Set(reflect.Append(refIndex, value))
		// case int
		default:
			f, err := strconv.Atoi(scan)
			if err != nil {
				return err
			}
			value := reflect.ValueOf(f)
			refIndex := ref.Field(index)
			refIndex.Set(reflect.Append(refIndex, value))
		}
	}
	return nil
}

/*
type T struct {
	Age int
	Name string
	Children []int
}
t := T{12, "someone-life", nil}
s := reflect.ValueOf(&t).Elem()

s.Field(0).SetInt(123) // 内置常用类型的设值方法
sliceValue := reflect.ValueOf([]int{1, 2, 3}) // 这里将slice转成reflect.Value类型
s.FieldByName("Children").Set(sliceValue)

作者：知乎用户
链接：https://www.zhihu.com/question/22783511/answer/24960616
来源：知乎
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
*/
