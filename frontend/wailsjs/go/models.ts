export namespace pb {
	
	export class Item {
	    id?: number;
	    name?: string;
	    description?: string;
	    type?: string;
	    icon?: string;
	    icon_large?: string;
	    members?: boolean;
	    current_price?: number;
	    current_trend?: string;
	    today_price_change?: number;
	    today_trend?: string;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.type = source["type"];
	        this.icon = source["icon"];
	        this.icon_large = source["icon_large"];
	        this.members = source["members"];
	        this.current_price = source["current_price"];
	        this.current_trend = source["current_trend"];
	        this.today_price_change = source["today_price_change"];
	        this.today_trend = source["today_trend"];
	    }
	}

}

