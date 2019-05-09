import { PipeTransform, Pipe } from '@angular/core';

@Pipe({
  name: 'callback',
  pure: true
})
export class CallbackPipe implements PipeTransform {
  transform(items: any[], filter: string) {

    if (!items || !filter) {
      return items;
    }

    let filteredItems = [];

    for(let i of items){
    	let itemString = String(i).toLowerCase();
    	if(itemString.includes(filter.toLowerCase())){
    		filteredItems.push(i);
    	}
    }

    return filteredItems;
  }
}
