const app = window.app;
const store = window.store;
const t = window.t;
const watch = window.watch;

const detailLayout = (title, data) => {
  const watchers = [];

  function selectRecord(record) {
    data.activeRecordId = record?.id || "";
  }

  return t.div(
    {
      className: "page detail-layout",
      style: "display: flex; flex-direction: column;",
      onunmount: () => watchers.forEach((watcher) => watcher?.unwatch()),
    },
    t.div(
      {
        style: "display: flex; gap: 20px; height: 100%;",
      },
      app.components.pageSidebar(
        {
          style: "margin: 20px;",
          onmount: (element) => {
            watchers.push(
              watch(
                () => data.activeRecordId,
                async () => {
                  await new Promise((resolve) => setTimeout(resolve, 0));

                  const activeNavItem =
                    element?.querySelector(".nav-item.active");
                  const details = activeNavItem?.closest("details");
                  if (details) {
                    details.open = true;
                    activeNavItem?.scrollIntoView({ block: "nearest" });
                  }
                },
              ),
            );
          },
        },
        t.h1(
          {
            style:
              "padding: 20px; border-bottom: 1px solid var(--surfaceAlt2Color);",
          },
          title,
        ),
        () =>
          t.nav(
            { className: "sidebar-content scrollable" },
            t.details(
              { className: "nav-group", open: true },
              t.summary(
                {
                  tabIndex: -1,
                  onfocusout: () => false,
                  onclick: () => false,
                  onkeyup: () => false,
                },
                "Today",
              ),
              () => {
                const presentableFields = data.presentableFields;

                return data.records.map((item) =>
                  t.button(
                    {
                      "html-data-record-id": () => item.id,
                      type: "button",
                      className: () =>
                        `nav-item responsive-close ${item.id == data.activeRecordId ? "active" : ""}`,
                      style:
                        "display: flex; flex-direction: column; align-items: flex-start; gap: 0",
                      onclick: (event) => {
                        event.preventDefault();
                        selectRecord(item);
                      },
                    },
                    t.p(
                      {
                        style:
                          "margin: 0; font-size: 14px; color: var(--surfaceTxtColor);",
                      },
                      presentableFields
                        .map((field) => item[field.name])
                        .filter((value) => value != null && value !== "")
                        .join(" · ") ||
                        item.title ||
                        item.id,
                    ),
                    t.p(
                      {
                        style: "margin: 0; font-size: 12px;",
                      },
                      item.created,
                    ),
                  ),
                );
              },
            ),
          ),
      ),
      t.div(
        {
          style: "flex: 1; padding: 20px;",
        },
        () => {
          const activeRecord = data.records.find(
            (item) => item.id == data.activeRecordId,
          );
          if (!activeRecord) {
            return "Select a record from the sidebar to view details.";
          }

          return t.pre(null, JSON.stringify(activeRecord, null, 2));
        },
      ),
    ),
  );
};

app.routes.superuserOnly("#/tools/{collectionName}", (route) => {
  const collectionName = route.params.collectionName;

  const data = store({
    collection: {},
    records: [],
    activeRecordId: "",
    get presentableFields() {
      if (!data.collection?.id || !Array.isArray(data.collection.fields)) {
        return [];
      }

      const result = data.collection.fields
        .filter((field) => field.presentable)
        .sort((field1, field2) => {
          const priority1 = app.fieldTypes[field1.type]?.summaryPriority ?? 0;
          const priority2 = app.fieldTypes[field2.type]?.summaryPriority ?? 0;
          return priority1 - priority2;
        });

      if (!result.length) {
        for (const name of app.utils.fallbackPresentableProps) {
          const field = data.collection.fields.find(
            (item) => item.name === name,
          );
          if (field) {
            result.push(field);
            break;
          }
        }
      }

      return result;
    },
  });

  async function loadRecords() {
    data.collection = await app.pb.collections.getOne(collectionName);
    data.records = await app.pb.collection(collectionName).getFullList();
  }
  loadRecords();

  return detailLayout(collectionName, data);
});

const tools = app.pb.collection("_tools").getFullList();
tools.then((tools) => {
  tools.forEach((tool) => {
    app.store.headerLinks.unshift({
      href: "#/tools/" + tool.targetCollectionName,
      icon: tool.icon ?? "",
      label: tool.targetCollectionName,
    });
  });
});

app.collectionTypes.base.tabs["Tool UI"] = function (upsertData) {
  return t.div(null, "extra fields...");
};
