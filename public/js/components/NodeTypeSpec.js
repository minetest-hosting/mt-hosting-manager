export default {
    props: ["nodetype"],
    template: /*html*/`
    <span>
        <i class="fa-solid fa-microchip" :title="nodetype.cpu_count + ' Processor(s)'"></i>
        {{nodetype.cpu_count}}
        <i class="fa-solid fa-memory" :title="nodetype.ram_gb + ' gigabyte RAM'"></i>
        {{nodetype.ram_gb}} GB
        <i class="fa-solid fa-hard-drive" :title="nodetype.disk_gb + ' gigabyte Disk'"></i>
        {{nodetype.disk_gb}} GB
        <i class="fa-solid fa-bolt" v-if="nodetype.dedicated" title="Dedicated machine"></i>
    </span>
    `
};