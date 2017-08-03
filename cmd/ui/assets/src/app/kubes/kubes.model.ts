export class Kube {
  public name: string;
  public id: number;
  public masterNodeSize: string;
  public status: boolean;

  constructor(name: string, desc: string) {
    this.name = name;
    this.masterNodeSize = desc;
    this.status = false;
  }
}
