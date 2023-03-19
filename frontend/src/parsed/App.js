import Tree from "react-d3-tree";
import { JsonTable } from "react-json-to-html";
import "./App.css";
import { Component } from "react";
import axios from "axios";
import { useLocation } from "react-router-dom";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      name: "React",
      showHideReturn: false,
    };

    this.setJson = this.setJson.bind(this);
    this.resetJson = this.resetJson.bind(this);
    this.siteStruct = this.props.struct;
    this.hierarchy = this.props.hierarchy;
    this.json = this.siteStruct;
  }

  setJson(node) {
    axios
      .get("http://localhost:1140/details?link=" + node.data.name)
      .then((response) => {
        const { text, ...data } = response.data;
        if (typeof data.metadata != "string") {
          data.metadata = JSON.stringify(data.metadata);
        }

        data.text = text;
        this.json = data;
        this.setState({ showHideReturn: true });
        this.forceUpdate();
      })
      .catch((err) => {
        console.log(err);
      });
  }

  resetJson() {
    this.json = this.siteStruct;
    this.setState({ showHideReturn: false });
    this.forceUpdate();
  }

  render() {
    return (
      <div className="App">
        <div className="Tree">
          <Tree
            data={this.hierarchy}
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
        <p>{JSON.stringify(this.state)}</p>
      </div>
    );
  }
}

function prepareJson(url, only_this_page, force, exclude) {
  console.log(url, only_this_page, force, exclude);
  // axios
  //   .get("http://localhost:1140/parse_site?url=" + "")
  //   .then((response) => {
  // const { text, ...data } = response.data;
  //   })
  //   .catch((err) => {
  //     console.log(err);
  //   });

  // const { hierarchy, ...rest } = data;
  return [
    { name: "test", childs: null },
    { name: "test", childs: null },
  ];
}

const WrappedComponent = (props) => {
  const location = useLocation();
  const { url, only_this_page, force, exclude } = location.state;

  let res = prepareJson(url, only_this_page, force, exclude);
  const struct = res[0],
    hierarchy = res[1];

  return <App struct={struct} hierarchy={hierarchy} {...props} />;
};

export default WrappedComponent;
