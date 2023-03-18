import Tree from "react-d3-tree";
import { JsonTable } from "react-json-to-html";
import { data, details } from "./mock";
import "./App.css";
import { Component } from "react";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: "React",
      showHideReturn: false,
    };

    this.setJson = this.setJson.bind(this);
    this.prepareJson = this.prepareJson.bind(this);
    this.prepareTree = this.prepareTree.bind(this);
    this.resetJson = this.resetJson.bind(this);
    this.siteStruct = this.prepareJson();
    this.json = this.siteStruct;
  }

  prepareJson() {
    const { hierarchy, ...rest } = data;
    return rest;
  }

  setJson() {
    if (typeof details.metadata != "string") {
      details.metadata = JSON.stringify(details.metadata);
    }

    this.json = details;
    this.setState({ showHideReturn: true });
    this.forceUpdate();
  }

  resetJson() {
    this.json = this.siteStruct;
    this.setState({ showHideReturn: false });
    this.forceUpdate();
  }

  prepareTree() {
    return data.hierarchy;
  }

  render() {
    return (
      <div className="App">
        <div className="Tree">
          <Tree
            data={this.prepareTree()}
            initialDepth={0}
            collapsible={true}
            transitionDuration={0}
            nodeSize={{ x: 400, y: 140 }}
            onNodeClick={this.setJson}
          />
        </div>
        <div className="Json">
          {this.state.showHideReturn && (
            <button id="return">
              <img
                src="https://icon-library.com/images/return-icon-png/return-icon-png-3.jpg"
                alt="return"
                width="30"
                onClick={this.resetJson}
              />
            </button>
          )}
          <JsonTable json={this.json} />
        </div>
      </div>
    );
  }
}

export default App;
