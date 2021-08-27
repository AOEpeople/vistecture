/*
    File contains methods to manipulate or read the view
 */

import $ from "jquery";
export default class layout {}

layout.GetGraphConfiguration = function () {
  return {
    hierarchicalSortMethod: $("#networkHierarchicalSortMethod").val(),
    nodeStyle: $("#networkNodeStyle").val(),
    layout: $("#networkLayout").val(),
    physics: $("#networkPhysics").prop("checked"),
    clusterGroups: $("#networkClusterGroups").val(),
    filterGroups: $("#networkFilterGroups").val(),
  };
};

layout.SetGraphPresets = function (preset) {
  let config = {};
  switch (preset) {
    case "hieradefault":
      config = {
        hierarchicalSortMethod: "directed",
        nodeStyle: "detailed",
        layout: "hierarchical",
        physics: false,
      };
      break;
    case "netdefault":
      config = {
        hierarchicalSortMethod: "",
        nodeStyle: "small",
        layout: "network",
        physics: true,
      };
      break;
  }

  $("#networkHierarchicalSortMethod").val(config["hierarchicalSortMethod"]);
  $("#networkLayout").val(config["layout"]);
  $("#networkNodeStyle").val(config["nodeStyle"]);
  $("#networkPhysics").prop("checked", config["physics"]);
};

layout.ShowSideContentModal = function (title, commontab, servicetab, deptab) {
  $("#sidecontent .close").click(function () {
    layout.hideSideContentModal();
  });
  $("#sidecontent #sidecontent-title").html(title);
  $("#sidecontent #common").html(commontab);
  $("#sidecontent #services").html(servicetab);
  $("#sidecontent #dependencies").html(deptab);
  $("#sidecontent").fadeIn();
};

layout.SetDocumentsMenu = function (documents) {
  for (var i in documents) {
    let document = documents[i];
    let urlname = encodeURI(document);
    $("#documentations").append(
      `<a class="dropdown-item" href="/documents/${urlname}" target="blank">${document}</a>`
    );
  }
};

layout.hideSideContentModal = function () {
  $("#sidecontent").fadeOut();
};
