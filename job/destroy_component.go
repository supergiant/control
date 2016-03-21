package job

// type DestroyComponentMessage struct {
// 	AppName       string
// 	ComponentName string
// }
//
// // DestroyComponent implements job.Performable interface
// type DestroyComponent struct {
// 	db   *storage.Client
// 	kube *guber.Client
// }
//
// func (j DestroyComponent) MaxAttempts() int {
// 	return 20
// }
//
// func (j DestroyComponent) Perform(data []byte) error {
//   message := new(DestroyComponentMessage)
//   if err := json.Unmarshal(data, message); err != nil {
//     return err
//   }
//
//   component, err := j.db.ComponentStorage.Get(message.AppName, message.ComponentName)
//   if err != nil {
//     return err
//   }
//
//   // Delete services
//   for _, service := range component.Services() {
//     if err = kube.Services(message.AppName).Delete(service.Metadata.Name); err != nil {
//       return err
//     }
//   }
//
//   // Delete replication controllers
//   // Delete pods
//   for _, instance := range component.ActiveDeployment().Instances() {
//     if rc := instance.ReplicationController(); rc != nil {
//       if err = kube.ReplicationControllers(message.AppName).Delete(rc.Metadata.Name); err != nil {
//         return err
//       }
//     }
//
//     if pod := instance.Pod(); pod != nil {
//       if err = kube.Pods(message.AppName).Delete(pod.Metadata.Name); err != nil {
//         return err
//       }
//     }
//   }
//
//   // Delete volumes
//   for _, volume := range component.Volumes() {
//     if err = volume.Destroy(); err != nil {
//       return err
//     }
//   }
//
//
//
//   // // Delete volumes
//   // volMgr := &cloud.VolumeManager{}
//   // for _, volDef := range component.Volumes {
//   //   volName := fmt.Sprintf("%s-%d", volDef.Name, instanceID)
//   //   volume, err := volMgr.Find(volName)
//   //   if err != nil {
//   //     continue
//   //   }
//   //
//   //   err = volMgr.WaitForAvailable(*volume.VolumeId)
//   //   if err != nil {
//   //     return err
//   //   }
//   //
//   //   err = volMgr.Delete(volName)
//   //   if err != nil {
//   //     return err
//   //   }
//   // }
// }
