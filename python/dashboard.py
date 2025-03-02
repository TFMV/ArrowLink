import streamlit as st
import grpc
import pyarrow as pa
import pyarrow.ipc as ipc
import pandas as pd
import plotly.express as px
import time
from proto.dataexchange_pb2_grpc import ArrowDataServiceStub
from proto.dataexchange_pb2 import Empty

st.set_page_config(page_title="ArrowLink Dashboard", layout="wide")

st.title("ArrowLink Real-time Data Dashboard")
st.markdown("Demonstrating high-performance data exchange between Go and Python")

# Sidebar controls
st.sidebar.header("Controls")
refresh_interval = st.sidebar.slider("Refresh interval (seconds)", 1, 10, 3)
auto_refresh = st.sidebar.checkbox("Auto-refresh", True)

# Connection status
conn_status = st.sidebar.empty()


# Create a connection to the gRPC server
@st.cache_resource
def get_stub():
    channel = grpc.insecure_channel("localhost:50051")
    return ArrowDataServiceStub(channel)


# Function to fetch data
def fetch_data():
    try:
        stub = get_stub()
        response_stream = stub.GetArrowData(Empty())
        for response in response_stream:
            reader = ipc.RecordBatchStreamReader(pa.BufferReader(response.payload))
            table = reader.read_all()
            return table.to_pandas()
    except Exception as e:
        st.error(f"Error fetching data: {e}")
        return None


# Main dashboard
col1, col2 = st.columns(2)

# Metrics
metrics_container = st.container()

# Charts
chart1 = col1.empty()
chart2 = col1.empty()
chart3 = col2.empty()
chart4 = col2.empty()

# Data table
data_table = st.empty()


# Function to update dashboard
def update_dashboard():
    start_time = time.time()
    conn_status.info("Fetching data...")

    df = fetch_data()
    if df is not None:
        fetch_time = time.time() - start_time
        conn_status.success(
            f"Connected to ArrowLink server (fetched in {fetch_time:.2f}s)"
        )

        # Update metrics
        with metrics_container:
            cols = st.columns(4)
            cols[0].metric("Total Records", len(df))

            # Check if 'category' column exists
            if "category" in df.columns:
                cols[1].metric("Categories", df["category"].nunique())
            else:
                cols[1].metric("Categories", "N/A")

            # Check if 'value' column exists
            if "value" in df.columns:
                cols[2].metric("Avg Value", f"{df['value'].mean():.2f}")
            else:
                cols[2].metric("Avg Value", "N/A")

            # Check if 'is_valid' column exists
            if "is_valid" in df.columns:
                cols[3].metric("Valid Records %", f"{df['is_valid'].mean()*100:.1f}%")
            else:
                cols[3].metric("Valid Records %", "N/A")

        # Update charts - only if required columns exist
        if "timestamp" in df.columns and "value" in df.columns:
            chart1.plotly_chart(
                px.line(
                    df.sort_values("timestamp").head(100),
                    x="timestamp",
                    y="value",
                    title="Recent Values Over Time",
                ),
                use_container_width=True,
                key="time_series_chart",
            )
        else:
            chart1.empty()
            st.info("Timestamp or value data not available")

        if "value" in df.columns:
            chart2.plotly_chart(
                px.histogram(df, x="value", title="Value Distribution"),
                use_container_width=True,
                key="histogram_chart",
            )
        else:
            chart2.empty()
            st.info("Value data not available")

        if "category" in df.columns and "value" in df.columns:
            chart3.plotly_chart(
                px.box(df, x="category", y="value", title="Values by Category"),
                use_container_width=True,
                key="box_chart",
            )
        else:
            chart3.empty()
            st.info("Category or value data not available")

        if "category" in df.columns:
            chart4.plotly_chart(
                px.pie(df, names="category", title="Records by Category"),
                use_container_width=True,
                key="pie_chart",
            )
        else:
            chart4.empty()
            st.info("Category data not available")

        # Update data table
        data_table.dataframe(df.head(10))
    else:
        conn_status.error("Failed to connect to ArrowLink server")


# Initial update
update_dashboard()

# Auto-refresh implementation
if auto_refresh:
    # Use a safer approach for auto-refresh
    st.markdown("**Auto-refreshing data...**")
    time.sleep(refresh_interval)
    st.rerun()
