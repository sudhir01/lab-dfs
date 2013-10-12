package viewservice

type View struct {
  Viewnum     int
  Primary     string
  Backup      string
  PrimaryView int
  BackupView  int
}

func (view *View) HasPrimaryAck() bool {
    return (view.Viewnum == INITIAL_VIEW) || (view.PrimaryView == view.Viewnum)
}
