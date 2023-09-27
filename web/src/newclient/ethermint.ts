import { Message } from '@bufbuild/protobuf';
import {
  TxPayload,
  TxContext,
} from '@evmos/transactions'
import { createTxRaw } from '@evmos/proto'
import { createTransactionWithMultipleMessages } from '@evmos/proto';
import { createEIP712, generateFee, generateMessageWithMultipleTransactions, } from '@evmos/eip712';
import Long from 'long'
import {
  AccountResponse,
  BroadcastMode,
  TxToSend,
  generateEndpointAccount,
  generateEndpointBroadcast,
  generatePostBodyBroadcast,
} from '@evmos/provider'
import { fromBase64 } from '@cosmjs/encoding';
import { chainDescriptor } from '../keplr';
import { bech32 } from 'bech32';
import { ETH } from '@evmos/address-converter';
import { Buffer } from 'buffer';

export async function fetchAccount(
  address: string,
) {
  const queryEndpoint = new URL(generateEndpointAccount(address), chainDescriptor.rest);

  const restOptions = {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
  }

  // Note that the node will return a 400 status code if the account does not exist.
  const rawResult = await fetch(
    queryEndpoint,
    restOptions,
  );

  const result = await rawResult.json();

  // Note that the `pub_key` will be `null` if the address has not sent any transactions.
  return result as AccountResponse;
}

export function buildTransaction(
  context: TxContext,
  msgs: Message<any>[],
) {
  const wrappedMsgs = msgs.map((msg) => ({
    message: msg,
    path: msg.getType().typeName,
  }));

  const txRaw = createTransactionWithMultipleMessages(
    wrappedMsgs,
    "",
    context.fee.amount,
    context.fee.denom,
    parseInt(context.fee.gas, 10),
    'ethsecp256',
    context.sender.pubkey,
    context.sender.sequence,
    context.sender.accountNumber,
    context.chain.cosmosChainId,
  )
  const feeObject = generateFee(context.fee.amount, context.fee.denom, context.fee.gas, context.sender.accountAddress);
  const msg = generateMessageWithMultipleTransactions(
    context.sender.accountNumber.toString(),
    context.sender.sequence.toString(),
    context.chain.cosmosChainId,
    context.memo,
    feeObject,
    wrappedMsgs,
  );

  const tx: TxPayload = {
    signDirect: txRaw.signDirect,
    legacyAmino: txRaw.legacyAmino,
    eipToSign: createEIP712(wrappedMsgs, context.chain.chainId, msg),
  }

  return tx;
}

export async function signTransactionKeplr(
  context: TxContext,
  tx: TxPayload,
) {
  const { chain, sender } = context

  const { signDirect } = tx

  const signResponse = await window?.keplr?.signDirect(
    chain.cosmosChainId,
    sender.accountAddress,
    {
      bodyBytes: signDirect.body.toBinary(),
      authInfoBytes: signDirect.authInfo.toBinary(),
      chainId: chain.cosmosChainId,
      accountNumber: new Long(sender.accountNumber),
    },
  )

  if (!signResponse) {
    throw new Error('No sign response');
  }

  const signatures = [
    fromBase64(signResponse.signature.signature),
  ]

  const { signed } = signResponse

  const signedTx = createTxRaw(
    signed.bodyBytes,
    signed.authInfoBytes,
    signatures,
  )

  return signedTx;
}

export const ethToFusion = (ethAddress) => {
    const data = ETH.decoder(ethAddress);
    return bech32.encode('qredo', bech32.toWords(data));
};
const fusionToEth = (fusionAddress) => {
    const { prefix, words } = bech32.decode(fusionAddress);
    return ETH.encoder(Buffer.from(bech32.fromWords(words)));
};

export async function signTransactionMetamask(
  context: TxContext,
  tx: TxPayload,
) {
  const { sender } = context

  const senderHexAddress = fusionToEth(sender.accountAddress)
  const eip712Payload = JSON.stringify(tx.eipToSign) // TODO: Fix this malformed payload
  console.log(eip712Payload)

  const signature = await window?.ethereum?.request({
    method: 'eth_signTypedData_v4',
    params: [senderHexAddress, eip712Payload],
  })
  console.log(signature)

  const signatureBytes = Buffer.from(signature.replace('0x', ''), 'hex')

  const { signDirect } = tx
  const bodyBytes = signDirect.body.toBinary()
  const authInfoBytes = signDirect.authInfo.toBinary()

  const signedTx = createTxRaw(
    bodyBytes,
    authInfoBytes,
    [signatureBytes],
  )

  return signedTx;
}

export async function broadcastTransaction(
  signedTx: TxToSend,
  broadcastMode: BroadcastMode | undefined = BroadcastMode.Sync,
) {
  const postOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: generatePostBodyBroadcast(signedTx, broadcastMode),
  }

  const broadcastEndpoint = new URL(generateEndpointBroadcast(), chainDescriptor.rest);
  const broadcastPost = await fetch(
    broadcastEndpoint,
    postOptions,
  )

  const { tx_response } = await broadcastPost.json()
  if (tx_response.code) {
    throw new Error("Error from chain node: " + tx_response.raw_log)
  }

  return tx_response as TxResponse;
}

export interface TxResponse {
  code: number,
  codespace: string,
  data: string,
  events: any[],
  gas_used: string,
  gas_wanted: string,
  height: string,
  info: string,
  logs: any[],
  raw_log: string,
  timestamp: string,
  tx: null,
  txhash: string,
}

