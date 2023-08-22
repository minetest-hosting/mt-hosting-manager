
export default {
	template: /*html*/`
    <div class="row">
        <div class="col-md-2"></div>
        <div class="col-md-8">
            <div class="card">
                <div class="card-header">
                    <slot name="title"></slot>
                </div>
                <div class="card-body">
                    <slot></slot>
                </div>
            </div>
        </div>
    </div>
	`
};
