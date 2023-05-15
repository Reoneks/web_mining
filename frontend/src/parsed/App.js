import Tree from "react-d3-tree";
import { JsonToTable } from "react-json-to-table";
import "./App.css";
import { Component } from "react";
import axios from "axios";
import { useLocation, Link } from "react-router-dom";
import _ from "lodash";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      showHideReturn: false,
      siteStruct: { Loading: "wait..." },
      hierarchy: [{ name: "Loading..." }],
    };

    this.setJson = this.setJson.bind(this);
    this.resetJson = this.resetJson.bind(this);
    this.url = this.props.url;
    this.only_this_page = this.props.only_this_page;
    this.force = this.props.force;
    this.exclude = this.props.exclude;
  }

  async getData(url, only_this_page, force, exclude) {
    const searchParams = new URLSearchParams();
    searchParams.append("url", url);
    searchParams.append(
      "only_this_page",
      only_this_page === undefined ? false : only_this_page
    );
    searchParams.append("force_collect", force === undefined ? false : force);
    exclude.forEach((excl,i) => {
      searchParams.append("exclude", excl);
    });

    const res = await axios(
      "http://localhost:1140/parse_site?" + searchParams.toString()
    );
    return res.data;
  }

  componentDidMount() {
    if (_.isEqual(this.state.siteStruct, { Loading: "wait..." })) {
      this.getData(this.url, this.only_this_page, this.force, this.exclude)
        .then((data) => {
          const { hierarchy, ...rest } = data;
          this.json = rest;
          this.setState({ siteStruct: rest, hierarchy: hierarchy });
          this.forceUpdate();
        })
        .catch((err) => {
          console.log(err);
        });
    }
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
    this.json = this.state.siteStruct;
    this.setState({ showHideReturn: false });
    this.forceUpdate();
  }

  render() {
    return (
      <div className="App">
        <div className="Tree">
          <Link to="/">
            <button id="close">
              <img
                src="https://cdn-icons-png.flaticon.com/512/463/463612.png"
                alt="close"
                width="30"
              />
            </button>
          </Link>
          <Tree
            data={this.state.hierarchy}
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
          <JsonToTable json={this.json} />
        </div>
      </div>
    );
  }
}

const WrappedComponent = (props) => {
  const location = useLocation();
  const { url, only_this_page, force, exclude } = location.state;

  return (
    <App
      url={url}
      only_this_page={only_this_page}
      force={force}
      exclude={exclude}
      {...props}
    />
  );
};

export default WrappedComponent;
