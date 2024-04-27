import React from "react";
import Graph from "react-graph-vis";
import uuid from "react-uuid";
import "./Graph.css"

const GraphShow = ({ node_list, edge_list }) => {
    // console.log(node_list)
    // console.log(edge_list)

    const graph = {
        nodes: node_list,
        edges: edge_list
    };

    const hierarchiallayout = {
        // hierarchicalLayout: {
        enabled: true,
        levelSeparation: 200,
        nodeSpacing: 300
        // },
    }

    const options = {
        autoResize: false,
        width: "750px",
        height: "500px",
        nodes: {
            shape: "star"
        },
        layout: {
            improvedLayout: false,
            hierarchical: {
                enabled: true,
                levelSeparation: 150,
                nodeSpacing: 500,
                treeSpacing: 200,
            }
        },
        physics: {
            enabled: false
        }
    }

    return (
        <div className="graph-container">
            <Graph
                key={uuid()}
                graph={graph}
                options={options}
                hierarchiallayout={hierarchiallayout}
            />
        </div>
    )
}

export default GraphShow;