/*
 File contains helper to work with vis network
 */

export default class visNetworkHelper {}

visNetworkHelper.InitNetwork = function (network, clickCallBack) {
  network.on("click", function (params) {
    clickCallBack(params);
  });
  visNetworkHelper.prepareNetworkForClustering(network);
};

visNetworkHelper.prepareNetworkForClustering = function (network) {
  network.on("doubleClick", function (params) {
    if (params.nodes.length == 1) {
      if (network.isCluster(params.nodes[0]) == true) {
        network.openCluster(params.nodes[0]);
      }
    }
  });
};

visNetworkHelper.ClusterByApplicationGroups = function (
  network,
  groups,
  clusterNodeRenderer
) {
  let groupNamesByLevel = [];
  for (var i in groups) {
    let groupName = groups[i];
    let groupPath = groupName.split("/");
    if (typeof groupNamesByLevel[groupPath.length] == "undefined") {
      groupNamesByLevel[groupPath.length] = {};
    }
    groupNamesByLevel[groupPath.length][groupName] = true;
  }
  // Cluster but start from deepes level for Clusters in Clusters
  for (let i = groupNamesByLevel.length - 1; i > 0; i--) {
    let groupNamesObject = groupNamesByLevel[i];
    if (typeof groupNamesObject == "undefined") {
      continue;
    }
    for (key in groupNamesObject) {
      let useKey = key;
      visNetworkHelper.clusterByApplicationGroup(
        network,
        useKey,
        clusterNodeRenderer
      );
    }
  }
};

//_clusterByApplicationGroup - does the clustering in vis network and clusters all nodes with groupName as node title
visNetworkHelper.clusterByApplicationGroup = function (
  network,
  groupName,
  clusterNodeRenderer
) {
  let groupPath = groupName.split("/");
  let parentGroupName = "";
  if (groupPath.length > 1) {
    parentGroupName = groupPath.slice(0, groupPath.length - 1).join("/");
  }

  clusterNode = clusterNodeRenderer(groupName);
  if (typeof clusterNode != "object") {
    console.error("Callback clusterNodeRenderer returned no object!");
    return;
  }
  //assign also appGroup - to allow for further clustering (cluster in cluster)
  clusterNode.appGroup = parentGroupName;

  network.cluster({
    joinCondition: function (childOptions) {
      if (childOptions.appGroup) {
        return (
          childOptions.appGroup.toLowerCase().trim() ==
          groupName.toLowerCase().trim()
        );
      }
      return false;
    },
    clusterNodeProperties: clusterNode,
  });
};
