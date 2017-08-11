import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';

@Injectable()
export class HeaderService {
    HeaderResponse = new Subject<any>();

    constructor() {}
}
