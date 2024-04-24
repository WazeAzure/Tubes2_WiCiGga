import React from "react";
import Graph from "react-graph-vis";
import uuid from "react-uuid";

const GraphShow = () => {
    
    const colorlist = ["#ce76fe", "#ed6f71", "#feae33", "#5358e2", "#fec46b"];

    const graph = {
        nodes: [
            { id: 1, label: "Node 1", title: "node 1 tootip text", level: -1 },
            { id: 2, label: "Node 2", title: "node 2 tootip text", level: -2 },
            { id: 3, label: "Node 3", title: "node 3 tootip text", level: -2 },
            { id: 4, label: "Node 4", title: "node 4 tootip text", level: -2 },
            { id: 5, label: "Node 5", title: "node 5 tootip text", level: -3 }
        ],
        edges: [
            { from: 1, to: 2 },
            { from: 1, to: 3 },
            { from: 1, to: 4 },
            { from: 2, to: 5 },
            { from: 3, to: 5 },
            { from: 4, to: 5 }
        ]
    };

    const options = {
        width: "750px",
        height: "500px",
        clustering: true,
        hierarchicalLayout: {
            enabled: true,
            levelSeparation: 200,
            nodeSpacing: 300
        },
        nodes: {
            radius: 100,
            shape: "circle"
        },
        physics: {
            barnesHut: {
                enabled: true,
                gravitationalConstant: -2000,
                centralGravity: 0.1,
                springLength: 100,
                springConstant: 0.05,
                damping: 0.09
            },
            repulsion: {
                centralGravity: 0.1,
                springLength: 50,
                springConstant: 0.05,
                nodeDistance: 100,
                damping: 0.09
            },
        },

    }

    return (
        <div className="graph-container" style={{backgroundColor: "#F7FFF7", border:"1px solid red"}}>
            <Graph
            key={uuid()}
            graph={graph}
            options={options}
            />
        </div>
    )
}

export default GraphShow;