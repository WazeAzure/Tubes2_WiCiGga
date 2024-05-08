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
        autoResize: true,
        width: "750px",
        height: "500px",
        nodes: {
            shape: "star"
        },
        layout: {
            improvedLayout: false,
            hierarchical: {
                enabled: true,
                levelSeparation: 200, // Adjust the vertical separation between levels
              },
        },
        physics: {
            enabled: true,
            repulsion: {
              nodeDistance: 200, // Increase this value to reduce repulsion between nodes
            },
            hierarchicalRepulsion: {
              centralGravity: 0, // Set central gravity to 0 to encourage nodes with the same level to be closer
              springLength: 1, // Adjust the spring length between nodes
              springConstant: 0.05, // Adjust the spring constant
              nodeDistance: 200, // Adjust the node distance
              damping: 0.9, // Adjust damping
            },
          },
    }

    const clustering = {
        enabled: true
    }

    return (
        <div className="graph-container">
            <Graph
                key={uuid()}
                graph={graph}
                options={options}
                hierarchiallayout={hierarchiallayout}
                clustering={clustering}
            />
        </div>
    )
}

export default GraphShow;