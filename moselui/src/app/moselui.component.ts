import {NodeService} from "./node.service";
import {Component} from "@angular/core";
import {CHART_DIRECTIVES} from "angular2-highcharts/index";
import {MD_CARD_DIRECTIVES} from "@angular2-material/card/card";
import {MD_GRID_LIST_DIRECTIVES} from "@angular2-material/grid-list/grid-list";

@Component({
  moduleId: module.id,
  selector: 'moselui-app',
  templateUrl: 'moselui.component.html',
  styleUrls: ['moselui.component.css'],
  directives: [
    CHART_DIRECTIVES,
    MD_CARD_DIRECTIVES,
    MD_GRID_LIST_DIRECTIVES
  ],
  providers: [NodeService],
})
export class MoseluiAppComponent {

  nodes = [];

  options = {
    title: {text: 'simple chart'},
    series: [{
      data: [29.9, 71.5, 106.4, 129.2],
    }]
  };

  constructor(private nodeService:NodeService) {
    this.nodeService.getNodes()
      .subscribe(
        res => {
          console.log(res);
          this.nodes = res.Nodes
        },
        err => console.log('Shit happens: ' + err))
  }

}
