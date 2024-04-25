import React from "react";
import Graph from "react-graph-vis";
import uuid from "react-uuid";

const GraphShow = ({node_list, edge_list}) => {
    console.log(node_list)
    console.log(edge_list)

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
        }
    }

    return (
        <div className="graph-container" style={{backgroundColor: "#F7FFF7", border:"1px solid red"}}>
            <Graph
            key={uuid()}
            graph={graph}
            options={options}
            hierarchiallayout = {hierarchiallayout}
            />
        </div>
    )
}

export default GraphShow;